package usecase

import (
	"context"
)

type RAGUsecase interface {
	LoadDocs(ctx context.Context, path string) ([]string, error)
	SendChat(ctx context.Context, message string) (string, error)
}
