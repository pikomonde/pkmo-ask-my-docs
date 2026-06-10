---
title: "What is MCP? A Simple Guide I Wish I Had Earlier"
url: "https://blog.pikomo.top/what-is-mcp-a-simple-guide-i-wish-i-had-earlier/"
author: "Piko Monde"
saved: "2026-06-08T03:59:14.963Z"
source: "blog.pikomo.top"
---

# What is MCP? A Simple Guide I Wish I Had Earlier

“AI Agents”, “Agentic”, “OpenClaw”, and “MCP”. I thought it was only a _Buzz Word_. Until my friend dragged me into a discussion I’m not ready for.

I have watch several videos and read several articles about it. And I thought I understand it.

1.  **MCP**: How LLMs (like chatGPT, Gemini, Gemma, etc) talks to external tools (API, etc)
    
2.  **OpenClaw**: Personal bot that can manage your work (and the bot can have their own social media, 2026 is a weird year right? lol 😅 )
    

Then, my friend said: “_If you want to understand OpenClaw, you need to learn MCP first_”

I search again about MCP and OpenClaw. Then, It hits me. It all connected. So, OpenClaw (our personal “agentic” bot) is powered by MCP. But, **what is MCP (Model Context Protocol)**?

## How MCP (Model Context Protocol) Works (The Architecture)?

To simplify, let me give a really simple example:

Okay, let’s say there are 4 players:

1.  User,
2.  Service A,
3.  LLM,
4.  Tools (APIs to `/get-user-balance/:id`, `/get-user-monthly-expense/:id`).

And here is the steps how it works:

1.  User prompt “_I do owe my friend Barbara 50 USD, how much can I pay her this month?_”;
2.  Service A sent the user prompt + list of API to LLM;
3.  LLM response to ask both `/get-user-balance/:id` and `/get-user-monthly-expanse/:id` APIs;
4.  Service A curl to both APIs, get the info, and sent the data to LLM `{ "userId": 10001, "averageMonthlyExpense": 4800 }`and `{ "userId": 10001, "balance": 4900 }`;
5.  LLM response back “_In average you spent 4900 USD a month, your balance is 4800 USD, I bet you will owe 100 more to Barbara this month, lol._”

![Diagram of MCP architecture: User sends a prompt to MCP Client, which communicates with the LLM. MCP Client also connects to MCP Server, which calls two external APIs: /get-user-balance and /get-user-expenses](https://blog.pikomo.top/_astro/mcp-architecture-clean.ChU_79X9_DJwx4.svg) Click to zoom

Diagram of MCP architecture: User sends a prompt to MCP Client, which communicates with the LLM. MCP Client also connects to MCP Server, which calls two external APIs: `/get-user-balance` and `/get-user-expenses`.

**You might have a question:** “_Why not LLM directly connect to the tools?_”.

Because LLM is “text only” (technically, in multi modal, it can see and reply photo and video too). It can only spit text. It can’t do curl, or do API to other thing. So, it needs to “talk” to the Service A. Therefore, we call Service A as “MCP Client”.

**You might say:** “_You lied to me. I do chat via Claude, but I can see they are thinking, like searching via web, calculating, copy paste data, execute program._”.

Hey, hey, that’s because, currently, Gemini, ChatGPT, Claude, they are an MCP Client. Not a plain LLM only anymore. But, if you integrate them (prompt) using their API, you can see that they might not have that agentic ability. Even, they (the API) don’t store your chat history, so you need to manage by yourself (but it will be discussed in different topic).

**You might ask again:** “_How Service A (MCP Client) know the list of API’s? How we can register that the tools to MCP Client?_”.

Nice question! Because.. I lied. When I say “the tools”, it is actually another service that manage other APIs. And this “the tools” service is called MCP Server. So know, we are having 4 core items + the APIs.

**And, you probably said:** “_I was right, I know you lied 👿. But how the LLM know about the list of tools?_”.

You might see it in step 2, that Service A (MCP Client), asks the MCP Server (the tools that connected to other external service and API) to list down all descriptions. Then use that description and list of tools with user’s first prompt to the LLM.

## MCP Server and MCP Client (The Standards)

But what does that “**list of API**” actually look like? It is a JSON, something like this:

```
{
  "tools": [
    {
      "name": "get_user_balance",
      "description": "Get the current balance of a user",
      "parameters": {
        "user_id": {
          "type": "string",
          "description": "The ID of the user"
        }
      }
    },
    {
      "name": "get_user_monthly_expense",
      "description": "Get the average monthly expense of a user",
      "parameters": {
        "user_id": {
          "type": "string",
          "description": "The ID of the user"
        }
      }
    }
  ]
}
```

This JSON is the “text” that MCP Client (Service A) sents to the LLM: “_Hey, these are the tools you can use. Here’s what each one does, and here’s what you need to call them._” The LLM reads the description field, and it knows which tool to pick based on the user’s prompt.

Now, I hope now you also connect the dots on how MCP works. But again, this is oversimplification, and I also just learn about this concept. It is nice to share what I know.

* * *

More projects on [GitHub](https://github.com/pikomonde) · If this saved you some debugging time, [Ko-fi](https://ko-fi.com/pikomonde) is always appreciated.
