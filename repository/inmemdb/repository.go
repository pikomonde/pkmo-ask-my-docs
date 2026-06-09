package inmemdb

import (
	"pikomonde/ask-my-docs/cmd/config"
	"pikomonde/ask-my-docs/entity"
	"pikomonde/ask-my-docs/repository"
)

type InMemoryDBRepo struct {
	cfg *config.Config
	db  map[string]entity.Doc
}

func New(
	cfg *config.Config,
) repository.DBRepository {
	return &InMemoryDBRepo{
		cfg: cfg,
		db:  make(map[string]entity.Doc),
	}
}
