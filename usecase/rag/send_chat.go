package rag

import (
	"context"
)

func (u *RAGUsecase) SendChat(ctx context.Context, queryMessage string) (string, error) {

	// Step 1: First Round to LLM (asking for RAG embedding vectors)

	// Step 2: Search using RAG

	// Step 3: Second Round to LLM (asking for RAG Embedding vectors)

	return "", nil
}
