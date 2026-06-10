---
title: "How We Saved a High-Traffic IoT Service from 200 RPS to 20,000+ RPS (and a $42k+ AWS Bill)"
url: "https://blog.pikomo.top/how-we-saved-iot-service-200-to-20000-rps/"
author: "Piko Monde"
saved: "2026-06-08T04:11:43.303Z"
source: "blog.pikomo.top"
---

# How We Saved a High-Traffic IoT Service from 200 RPS to 20,000+ RPS (and a $42k+ AWS Bill)

In the world of high-concurrency systems, throwing more hardware at a problem is often the most expensive way to fail. Recently, I revisited some investigation logs and Go pprof profiles from a project I handled four years ago as a contractor for an Automobile IoT company. At the time, the company was managing telemetry for tens of thousands of connected vehicles.

The service was struggling with massive CPU utilization and scaling issues. Despite being backed by a significant cloud budget, the infrastructure was buckling under a load that, on paper, should have been manageable. This is a story of how we moved from a state of “throwing money at the fire” to a lean, high-performance architecture.

## The Infrastructure Bottleneck: Roughly 27 Nodes for 200 RPS

My first task at this company was to optimize our gateway server. Since I was a contractor and didn’t have direct access yet, the Engineering Lead and I sat down to review the dashboard together. When he showed me the infrastructure, I was floored.

The service was running on a cluster of roughly **27 large compute instances** on AWS. As a contractor I wasn’t shown the billing console directly, but back-of-the-envelope math told the story: at on-demand pricing for large compute-optimized instances (I’m using c6g.16xlarge for this calculation), this cluster was burning somewhere **north of $42,000 a month** — for this one service layer alone. Worse, CPU usage across all nodes was pegged at 90%+.

The Engineering Manager felt his hands were tied; his schedule was packed with meetings, and the rest of the team was primarily focused on DevOps. I was essentially the only one available and specifically hired to fix this mess. Looking at the code, the EM had a hunch: _“I think it’s the RBAC logic, especially when unmarshalling the user entity from Redis cache.”_

## Identifying the Bottleneck: A Lesson in Profiling

To understand why a cluster that large couldn’t handle a measly 200 RPS, I started my investigation using API call samples provided by one of the DevOps engineers.

### Setting Up Observability

I implemented the `net/http/pprof` package to capture live data under a simulated load, by adding a simple blank import:

```
import _ "net/http/pprof"
```

The service started exposing debugging endpoints at `/debug/pprof/`. This allowed me to capture live data while the service was under simulated load.

### The Local Benchmark Discrepancy

I began the investigation using `go-wrk`. I ran a standard 5-second test with 10 concurrent connections:

```
go-wrk -d 5 -c 10 http://localhost:8080/v1/telemetry
```

To my surprise, my local environment was hitting nearly **2,000 RPS** — ten times the production capacity. This discrepancy was confusing. I initially thought the bottleneck might be related to minor inefficiencies in the code that only compounded at scale.

### The First Round of Optimizations

At this point, I performed optimizations based on the initial pprof data.

**RBAC JSON Unmarshalling:** The service used a Redis-based Role-Based Access Control (RBAC) system. For every request, it fetched a massive `User` struct from Redis and unmarshalled the entire JSON object to check a single permission bit. I refactored this to store simple boolean flags or bitmasks in Redis. This reduced the CPU time spent on JSON decoding by nearly 30%.

The profile at this stage looked something like this — JSON unmarshalling is visible, eating around 20–30% of CPU, but nothing that would explain a full cluster meltdown:

![Go pprof flamegraph showing RBAC JSON unmarshalling consuming around 20-30% CPU — visible but not catastrophic](https://blog.pikomo.top/_astro/pprof-flamegraph-rbac-json-unmarshal-go.CoB0iBAe_zAC9Q.webp) Click to zoom

Figure 1A: Initial pprof profile — RBAC JSON unmarshalling is the largest consumer, but the overall picture still looks manageable. (Reproduced on a test service for illustration; the original production profiles are proprietary.)

I thought I had found the solution. But when we deployed the RBAC fix to production, nothing changed. The CPU stayed at 100%, and the throughput stayed at 200 RPS.

## The Breakthrough: Production vs. Local Payloads

The “Eureka” moment happened during a 1-on-1 session with my Engineering Manager. We compared my local requests with actual production traffic and realized I was **missing a specific header**.

When I added that header and ran pprof again, the trace revealed something entirely different. The request wasn’t hitting the standard user path; it was being routed to a specific IoT Middleware. In this middleware, there was a hidden process: it wasn’t checking session or refresh tokens. Instead, it was **validating every single incoming telemetry token using Bcrypt**.

As soon as I updated my local `go-wrk` script to include the production-spec device headers, my local CPU instantly spiked to 100%, and my RPS plummeted to exactly what we saw in production: **200 RPS**.

This is what the new profile looked like. Same service, same endpoint — spot the difference:

![Go pprof flamegraph showing bcrypt.CompareHashAndPassword consuming 95%+ of total CPU after adding the IoT device header](https://blog.pikomo.top/_astro/pprof-flamegraph-bcrypt-hotpath-cpu-saturation.m3JVh-wS_1uy9qX.webp) Click to zoom

Figure 1B: Profile after adding the IoT device header. bcrypt.CompareHashAndPassword has consumed almost everything — there is almost nothing else left. (Reproduced on a test service for illustration.)

## The Bcrypt Implementation Issue: Self-Inflicted DDoS

The flame graph was unmistakable. A single function was consuming over 80% of the total CPU time: `bcrypt.CompareHashAndPassword`.

We discovered that for every single telemetry ping — sent by thousands of vehicles every few seconds — the system was executing a Bcrypt comparison.

From an architectural perspective, this is a catastrophic anti-pattern. Bcrypt is designed by cryptographers to be **slow**. It uses a computational “cost factor” to ensure that even with massive hardware, brute-forcing a password takes an eternity. It is meant for _login_ endpoints, not _telemetry_ endpoints.

**We were essentially DDoS-ing ourselves.**

![This is Fine meme — a dog sits calmly in a room that is entirely on fire](https://blog.pikomo.top/_astro/this-is-fine-meme-self-inflicted-ddos.CkXvXshd_Z2hRuV7.webp) Click to zoom

Production CPU at 100%. “Oh, we are DDoS-ing ourselves.” — _This is Fine_, KC Green, Gunshow comic (2013). Via [Know Your Meme](https://knowyourmeme.com/memes/this-is-fine).

While the intention was to have high security for every data point, the implementation was impractical. For high-frequency IoT data, the industry standard is to use Bcrypt only during the initial handshake to exchange credentials for a lightweight session token or a short-lived JWT (JSON Web Token). Once authenticated, subsequent telemetry should be validated using symmetric keys or token lookups, which are orders of magnitude faster than Bcrypt.

## The Solution: Caching and Singleflight

Since we couldn’t update the firmware of thousands of vehicles overnight, we had to build a server-side “shield”.

### Implementing the Cache

I consulted with the Security team first. The whole point of Bcrypt is to prevent brute-forcing, so I asked: _“If we cache the successful validation result, is it still secure?”_ They gave us the green light.

We introduced a caching layer. After a successful Bcrypt validation, we stored the **SHA-256 hash** of the device credentials in Redis with a **30-second TTL** (Time to Live). For any subsequent request within that window, the server would simply compare the SHA-256 hashes — a process that takes nanoseconds — instead of running the Bcrypt algorithm.

### Resolving the Thundering Herd with Singleflight

After deploying the cache, we saw a massive improvement, but every 30 seconds, the CPU would spike again. This was the **“Thundering Herd”** problem. When the cache expired, all concurrent requests for the same device would see a cache miss and simultaneously trigger a Bcrypt operation.

To resolve this, we implemented [`golang.org/x/sync/singleflight`](https://pkg.go.dev/golang.org/x/sync/singleflight). The logic was simple: for a given device ID, only one Bcrypt operation should be “in flight” at any given time.

```
// Only one goroutine per deviceID executes this at a time
v, err, _ := g.Do(deviceID, func() (interface{}, error) {
    return validateBcrypt(storedHash, providedPassword)
})
```

By using singleflight, if 1,000 requests for `Device-A` arrive at the same time during a cache miss, **one** request performs the Bcrypt check, while the other 999 wait for that single result.

> **Note on multi-instance behavior:** Singleflight is local to the instance. In our multi-node setup, we might still perform one Bcrypt operation per instance during a cache miss — but this was a negligible cost compared to the thousands of operations we were doing previously.

> **Note on security:** While this fixed the performance, the system remains somewhat vulnerable to a focused DDoS. If an attacker cycles through thousands of different invalid tokens, the CPU would still spike because each unique invalid token triggers a new Bcrypt operation.

## Results: 90%+ Reduction Across the Board

The impact was immediate. During a final 1-on-1, my Engineering Manager decided to test the limits. He scaled the environment down to just two small/medium instances.

Here is what changed under the hood — from every request triggering bcrypt directly, to a layered shield where 99% of requests never touch it at all:

![Architecture diagram showing request flow before and after the optimization: before, every request hits bcrypt directly; after, a Redis cache and singleflight barrier absorb almost all load](https://blog.pikomo.top/_astro/request-flow-before-after-bcrypt-cache-singleflight.RKDPSihc_2d81MY.svg) Click to zoom

Request flow before (left) and after (right) the cache + singleflight fix. Cache hits return in nanoseconds; bcrypt is now rare and controlled.

We watched the dashboard in silence. The system didn’t just survive — it thrived.

Metric

Before

After

Throughput

200 RPS

20,000+ RPS

CPU Utilization

90–100%

<50%

Instance count

~27 nodes

2 smaller instances

Est. monthly cost

North of $42k/mo\*\*

Drastically reduced

\*_Estimated from node count and on-demand pricing for the instance family. Exact billing figures were not shared with me as a contractor._

After that meeting, I received an email that honestly made me a bit emotional. The EM sent a company-wide announcement, crediting me for fixing the system’s core stability and drastically slashing the AWS bill.

## Key Takeaways

1.  **Profile early, profile often.** Don’t guess where the bottleneck is. Use `pprof` to see exactly where the CPU cycles are going.
    
2.  **Production parity matters.** My local tests failed initially because I wasn’t using production-grade headers. Always ensure your load tests mimic real-world traffic patterns — including headers and metadata.
    
3.  **Security vs. performance.** Security is paramount, but expensive algorithms like Bcrypt don’t belong in high-frequency “hot paths.” Use them for handshakes, not for every data packet.
    
4.  **Beware of the Thundering Herd.** Caching is not a silver bullet. When dealing with high concurrency, always consider what happens when the cache expires. Tools like `singleflight` are essential in your Go toolkit.
    
5.  **Infrastructure is not a substitute for optimization.** Scaling horizontally can hide bad code for a while, but eventually, the technical debt will become too expensive to ignore.
    

* * *

More projects on [GitHub](https://github.com/pikomonde) · If this saved you some debugging time, [Ko-fi](https://ko-fi.com/pikomonde) is always appreciated.
