package rag

import (
	"pikomonde/ask-my-docs/cmd/config"
	"pikomonde/ask-my-docs/repository"
	"pikomonde/ask-my-docs/usecase"
)

type RAGUsecase struct {
	cfg            *config.Config
	AIRpository    repository.AIRepository
	InMemoryDBRepo repository.DBRepository
}

func New(
	cfg *config.Config,
	AIRpository repository.AIRepository,
	InMemoryDBRepo repository.DBRepository,
) usecase.RAGUsecase {
	return &RAGUsecase{
		cfg:            cfg,
		AIRpository:    AIRpository,
		InMemoryDBRepo: InMemoryDBRepo,
	}
}
