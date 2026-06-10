package googleai

import (
	"context"
	"fmt"

	"google.golang.org/genai"
)

func (r *GoogleAIRepo) SendMessage(ctx context.Context, contents []*genai.Content) (string, error) {
	genaiResp, err := r.googleAICli.Cli.Models.GenerateContent(ctx, r.cfg.GoogleAI.LLMModel, contents, nil)
	if err != nil {
		return "", fmt.Errorf("error when get generate content from the API")
	}

	if len(genaiResp.Candidates) == 0 {
		return "", fmt.Errorf("error llm model returns no candidates")
	}

	for _, part := range genaiResp.Candidates[0].Content.Parts {
		if !part.Thought {
			return part.Text, nil
		}
	}

	return "", fmt.Errorf("error llm model returns no non-thought text")
}
