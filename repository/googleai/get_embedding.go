package googleai

import (
	"context"
	"fmt"

	"google.golang.org/genai"
)

func (r *GoogleAIRepo) GetEmbedding(ctx context.Context, text string) ([]float32, error) {
	contents := []*genai.Content{
		genai.NewContentFromText(text, genai.RoleUser),
	}

	result, err := r.googleAICli.Cli.Models.EmbedContent(ctx, r.cfg.GoogleAI.EmbeddingModel, contents, nil)
	if err != nil {
		return nil, fmt.Errorf("error when get embedding vectors from the API")
	}

	if len(result.Embeddings) == 0 {
		return nil, fmt.Errorf("no embedding vectors returned from the API")
	}

	return result.Embeddings[0].Values, nil
}
