# Ask My Docs

This is a CLI app that implements RAG with simple Google AI.

## How to run:

```
GOOGLE_API_KEY=xxxx go run cmd/cli/app.go --docs ./articles --chunk-mode 4
```

## Args

|Name|Args|Default Value|Description|Notes|
|-|-|-|-|-|
|Article Directory Path|`docs`|`./docs`|path to docs folder||
|Chunk Mode|`chunk-mode`|`4`: ChunkModeParagraphWordSmart|mode of chuking: (1:ChunkModeNoChunk, 2:ChunkModeWord, 3:ChunkModeParagraph, 4:ChunkModeParagraphWordSmart)|**ChunkModeNoChunk:** The article is not chunk at all; **ChunkModeWord:** The article is chunk by word (the arg chunk size is used as word); **ChunkModeParagraph:** The article is chunk by paragraph (the arg chunk size is used as paragraph); **ChunkModeParagraphWordSmart:** The article is chunk by word, but snapped to ceil paragraph (the arg chunk size is used as word)|
|Chunk Size|`chunk`|`400` if words; `5` if paragraph|chunk size (default: 5 for paragraph, 400 for words)||
|Chunk Overlap|`chunk-overlap`|`80` if words; `1` if paragraph|overlapping chunk size (default: 1 for paragraph, 80 for words)||
|Google AI Api Key|`api-key`|(empty)|google-ai API key|you can use env var `GOOGLE_API_KEY` instead|
|Google AI Embedding Model|`embedding-model`|`gemini-embedding-2-preview`|google-ai's embedding model used||
|Google AI LLM Model|`llm-model`|`gemma-4-31b-it`|google-ai's llm model used||

## Type of Docs
In can read any txt and md files. It is based to use in tandem with [Articon](https://www.pikomo.top/projects/articon-go) extension to "save" articles from websites and download it locally.

## How it works?

### Loading Steps:
1. Load text documents;
2. Chunking each documents;
3. Transform each text into vector (embedding via Embedding Model);
4. Save to local in memory DB.

### Chat Steps:
1. It retrieve query from user, for example: "*What is algorithm?*";
2. **Step A:** hit LLM Model to get 4-8 RAG query words, for example: "*definition and explanation of algorithms*";
3. **Step B1:** Transform it into vectors (embedding) via Embedding Model;
4. **Step B2:** Search in DB using cosineSimilarity, basically just a dot product to each database items;
5. **Step B3:** Sort by similairty;
6. **Step C:** hit LLM Model with user's query and chunk text to get the final answer.
7. Send back to user.

