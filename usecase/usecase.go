package usecase

import (
	"context"
	"pikomonde/ask-my-docs/entity"
)

type RAGUsecase interface {
	LoadDocs(ctx context.Context, path string) (entity.Docs, error)
	SendChat(ctx context.Context, message string) (string, error)
}
