package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	cliGoogleAI "pikomonde/ask-my-docs/client/googleai"
	"pikomonde/ask-my-docs/cmd/config"
	rGoogleAI "pikomonde/ask-my-docs/repository/googleai"
	rInMemDB "pikomonde/ask-my-docs/repository/inmemdb"
	"pikomonde/ask-my-docs/usecase"
	uRAG "pikomonde/ask-my-docs/usecase/rag"

	"github.com/sirupsen/logrus"
)

func main() {
	logrus.SetLevel(logrus.DebugLevel)

	cliHandler, err := NewCLIHandler()
	if err != nil {
		log.Fatalln("error when init CLI handler", err.Error())
	}

	// new config
	cfg := cliHandler.Config
	ctx := context.Background()

	// new clients
	googleAICLi, err := cliGoogleAI.New(ctx, cfg)
	if err != nil {
		log.Fatalln("error when getting new clients", err.Error())
	}

	// new repositories
	googleAIRepo := rGoogleAI.New(cfg, googleAICLi)
	inMemDBRepo := rInMemDB.New(cfg)

	// new usecases
	ragUC := uRAG.New(cfg, googleAIRepo, inMemDBRepo)

	// new handler
	cliHandler.SetUsecase(ragUC)
	cliHandler.Start(ctx)
}

type CliHandler struct {
	Config *config.Config
	uc     usecase.RAGUsecase
}

func NewCLIHandler() (*CliHandler, error) {
	cfg := &config.Config{}

	var chunkMode int

	flag.StringVar(&cfg.DocsDirPath, "docs", "./docs", "path to docs folder")
	flag.IntVar(&chunkMode, "chunk-mode", 4, "mode of chuking: (1:ChunkModeNoChunk, 2:ChunkModeWord, 3:ChunkModeParagraph, 4:ChunkModeParagraphWordSmart)")

	flag.IntVar(&cfg.ChunkSize, "chunk", -1, "chunk size (default: 5 for paragraph, 400 for words)")
	flag.IntVar(&cfg.ChunkOverlap, "chunk-overlap", -1, "overlapping chunk size (default: 1 for paragraph, 80 for words)")

	flag.StringVar(&cfg.GoogleAI.APIKey, "api-key", "", "google-ai API key")
	flag.StringVar(&cfg.GoogleAI.EmbeddingModel, "embedding-model", "gemini-embedding-2-preview", "google-ai's embedding model used")
	flag.StringVar(&cfg.GoogleAI.LLMModel, "llm-model", "gemma-4-31b-it", "google-ai's llm model used")

	flag.Parse()

	cfg.ChunkMode = config.ChunkMode(chunkMode)

	if cfg.ChunkSize == -1 {
		if cfg.ChunkMode == config.ChunkModeParagraph {
			cfg.ChunkSize = 5
		} else {
			cfg.ChunkSize = 400
		}
	}

	if cfg.ChunkOverlap == -1 {
		if cfg.ChunkMode == config.ChunkModeParagraph {
			cfg.ChunkOverlap = 1
		} else {
			cfg.ChunkOverlap = 80
		}
	}

	if cfg.GoogleAI.APIKey == "" {
		cfg.GoogleAI.APIKey = os.Getenv("GOOGLE_API_KEY")
	}

	if cfg.GoogleAI.APIKey == "" {
		logrus.Info("Input Google AI API Key: ")
		reader := bufio.NewReader(os.Stdin)
		apiKey, err := reader.ReadString('\n')
		if err != nil {
			return nil, fmt.Errorf("failed to read API key: %w", err)
		}
		cfg.GoogleAI.APIKey = strings.TrimSpace(apiKey)

	}

	return &CliHandler{Config: cfg}, nil
}

func (h *CliHandler) SetUsecase(uc usecase.RAGUsecase) {
	h.uc = uc
}

func (h *CliHandler) Start(ctx context.Context) {
	// 1. validate config
	if h.Config == nil {
		panic("config is not set")
	}
	if h.uc == nil {
		panic("usecase is not set")
	}

	// 2. loading docs
	docs, err := h.uc.LoadDocs(ctx, h.Config.DocsDirPath)
	if err != nil {
		log.Fatalln("error when LoadDocs", err.Error())
	}
	for _, doc := range docs {
		logrus.Debugf("Document %s, Chunk %d [%d words] is loaded\n", doc.Path, doc.ChunkNo, len(strings.Split(doc.Text, " ")))
	}
	logrus.Debugln("Documents Loaded.")
	logrus.Debugln()

	// 3. interactive chat
	logrus.Infoln("==========================")
	logrus.Infoln("======= Let's Chat =======")
	logrus.Infoln("==========================")

	for {
		logrus.Info("User: ")
		reader := bufio.NewReader(os.Stdin)
		msg, err := reader.ReadString('\n')
		if err != nil {
			log.Printf("failed to read user chat: %s\n", err.Error())
		}

		replyMsg, err := h.uc.SendChat(ctx, msg)
		if err != nil {
			log.Fatalln("error when SendChat", err.Error())
		}
		logrus.Infof("Bot: %s\n", replyMsg)
	}
}
