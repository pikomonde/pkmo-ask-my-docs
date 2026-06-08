package googleai

import (
	"context"
	"log"
	"pikomonde/ask-my-docs/cmd/config"

	"google.golang.org/genai"
)

type GoogleAIClient struct {
	Cli *genai.Client
}

func New(ctx context.Context, cfg *config.Config) (*GoogleAIClient, error) {
	cliConfig := &genai.ClientConfig{
		APIKey: cfg.GoogleAI.APIKey,
	}

	// Initialize the centralized client.
	client, err := genai.NewClient(ctx, cliConfig)
	if err != nil {
		log.Fatalf("Failed to initialize GenAI client: %v", err)
	}

	return &GoogleAIClient{
		Cli: client,
	}, nil
}
