package inmemdb

import (
	"pikomonde/ask-my-docs/cmd/config"
	"pikomonde/ask-my-docs/repository"
)

type InMemoryDBRepo struct {
	cfg *config.Config
	db  map[string][]float32
}

func New(
	cfg *config.Config,
) repository.DBRepository {
	return &InMemoryDBRepo{
		cfg: cfg,
		db:  make(map[string][]float32),
	}
}
