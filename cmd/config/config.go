package config

type ChunkMode int

const (
	// No chunk
	ChunkModeNoChunk = iota + 1

	// Chunk text by number of words
	//
	// It is using N number of words, where N is from ChunkSize
	ChunkModeWord

	// Chunk text by number of paragraphs
	//
	// It is using N number of paragraphs, where N is from ChunkSize
	ChunkModeParagraph

	// Chunk text by number of words, but snapped (floored) into paragraph
	//
	// It loops through each paragraph, count the number of words.
	// If it more than ChunkSize, it stops, and only takes N-1 paragraphs
	ChunkModeParagraphWordSmart

	// TODO:
	//
	// Read through the MD formattings, by see the headings
	// ChunkModeParagraphWordHeadingSmart
)

type Config struct {
	GoogleAI     GoogleAIConfig
	DocsDirPath  string
	ChunkSize    int // Default is 400 words or 5 parapgraphs
	ChunkOverlap int // Default is 80 words or 1 paragraph, rule of thumb is 10-20% of ChunkSize
	ChunkMode    ChunkMode
}

type GoogleAIConfig struct {
	APIKey         string
	EmbeddingModel string // default gemini-embedding-2-preview
	LLMModel       string // default gemma-4-31b-it
}
