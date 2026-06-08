package googleai

import (
	cliGoogleAI "pikomonde/ask-my-docs/client/googleai"
	"pikomonde/ask-my-docs/cmd/config"
	"pikomonde/ask-my-docs/repository"
)

type GoogleAIRepo struct {
	cfg         *config.Config
	googleAICli *cliGoogleAI.GoogleAIClient
}

func New(
	cfg *config.Config,
	googleAICli *cliGoogleAI.GoogleAIClient,
) repository.AIRepository {
	return &GoogleAIRepo{
		cfg:         cfg,
		googleAICli: googleAICli,
	}
}
