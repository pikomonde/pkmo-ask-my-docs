package rag

import (
	"context"
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
	"google.golang.org/genai"
)

const (
	RAG_SYSTEM_MESSAGE = `You are a search query extractor for an article research assistant.
Given a user's question and recent chat history, output ONLY a concise search query (3-10 words) 
that would help find relevant articles to answer the question.
Output just the search query, no explanation, no quotes.`

	RESPONSE_TO_HUMAN_SYSTEM_MESSAGE = `You are an AI research assistant that helps users explore their saved articles.
You answer questions based on the user's personal article collection.
Be conversational, accurate, and cite article titles when relevant.
If the provided articles don't contain enough information, say so honestly.
Keep responses focused and useful — avoid padding.
Avoid any styling, remember that this is a CLI (command line interface) app.`
)

func (u *RAGUsecase) SendChat(ctx context.Context, queryMessage string) (string, error) {

	// Step 1: First Round to LLM (asking for RAG embedding vectors)
	ragMessage, err := u.AIRpository.SendMessage(ctx, []*genai.Content{{
		Role: genai.RoleModel, Parts: []*genai.Part{{Text: RAG_SYSTEM_MESSAGE}},
	}, {
		Role: genai.RoleUser, Parts: []*genai.Part{{Text: queryMessage}},
	}})
	if err != nil {
		logrus.Errorf("Fail to send query message %s: %v\n", queryMessage, err)
		return "", err
	}
	logrus.Debugf("System: [decoding message: %s]\n", ragMessage)

	embeddedMsg, err := u.AIRpository.GetEmbedding(ctx, ragMessage)
	if err != nil {
		logrus.Errorf("Fail to get query embedding %s: %v\n", ragMessage, err)
		return "", err
	}

	// Step 2: Search using RAG
	logrus.Debugf("System: querying docs...\n")
	queryResultDocs, err := u.InMemoryDBRepo.Search(embeddedMsg)
	if err != nil {
		logrus.Errorf("Fail to search embedding %v\n", err)
		return "", err
	}

	// Step 3: Second Round to LLM (asking for RAG Embedding vectors)
	logrus.Debugf("System: getting %d docs:\n", len(queryResultDocs))
	contextDocs := ""
	for rank, queryResultDoc := range queryResultDocs {
		if rank > 3 {
			break
		}
		logrus.Debugf("%s...\n", fmt.Sprintf("System: title \"%s\" [%d] (score: %2f): %s",
			queryResultDoc.Path, queryResultDoc.ChunkNo, queryResultDoc.CosineSimilarity, strings.ReplaceAll(queryResultDoc.Text, "\n", ""))[:140])
		contextDocs += fmt.Sprintf(`
--- Result rank %d: %s [chunk no: %d] (score: %2f) ---
%s
`,
			rank+1, queryResultDoc.Path, queryResultDoc.ChunkNo, queryResultDoc.CosineSimilarity,
			queryResultDoc.Text)
	}
	finalQuery := fmt.Sprintf(`%s
Relevant articles from your collection:
%s
`,
		queryMessage, contextDocs)

	logrus.Debugf("System: processing response...\n")
	finalMessage, err := u.AIRpository.SendMessage(ctx, []*genai.Content{{
		Role: genai.RoleModel, Parts: []*genai.Part{{Text: RESPONSE_TO_HUMAN_SYSTEM_MESSAGE}},
	}, {
		Role: genai.RoleUser, Parts: []*genai.Part{{Text: finalQuery}},
	}})
	if err != nil {
		logrus.Errorf("Fail to send query message %s: %v\n", queryMessage, err)
		return "", err
	}

	return finalMessage, nil
}
