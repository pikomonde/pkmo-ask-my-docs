package rag

import (
	"context"
	"os"
	"path/filepath"
	"pikomonde/ask-my-docs/cmd/config"
	"pikomonde/ask-my-docs/entity"
	"strings"

	"github.com/sirupsen/logrus"
)

func (u *RAGUsecase) LoadDocs(ctx context.Context, directoryPath string) (entity.Docs, error) {
	docs := make(entity.Docs, 0)

	// 1. List files from folder
	files, err := os.ReadDir(directoryPath)
	if err != nil {
		return nil, err
	}

	// 2. Iterate by files
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		filePath := filepath.Join(directoryPath, file.Name())

		// 3. Read files
		textBytes, err := os.ReadFile(filePath)
		if err != nil {
			logrus.Errorf("Fail to read %s: %v\n", file.Name(), err)
			continue
		}

		// 4. Chunking
		chunks := chunkText(string(textBytes), u.cfg.ChunkMode, u.cfg.ChunkSize, u.cfg.ChunkOverlap)

		for chunkID, chunkTxt := range chunks {
			// 5. Get Embedding
			embeddedText, err := u.AIRpository.GetEmbedding(ctx, chunkTxt)
			if err != nil {
				logrus.Errorf("Fail to get embedding %s: %v\n", file.Name(), err)
				continue
			}

			// 6. Append text
			doc := entity.Doc{
				Path:       file.Name(),
				ChunkNo:    chunkID,
				Text:       chunkTxt,
				Embeddings: embeddedText,
			}
			docs = append(docs, doc)

			// 7. Save Emedding
			err = u.InMemoryDBRepo.SaveDoc(doc.Key(), doc)
			if err != nil {
				logrus.Errorf("Fail to save doc %s: %v\n", file.Name(), err)
				continue
			}
		}

	}

	return docs, nil
}

func chunkText(txt string, chunkMode config.ChunkMode, chunkSize, chunkOverlap int) []string {
	switch chunkMode {
	case config.ChunkModeNoChunk:
		return []string{txt}
	case config.ChunkModeWord:
		result := make([]string, 0)

		txtSplit := strings.Split(txt, " ")
		chunkGap := chunkSize - chunkOverlap
		// fmt.Println("---> len(txtSplit)", len(txtSplit))
		// fmt.Println("---> chunkGap", chunkGap)
		// fmt.Println("---> len(txtSplit)/chunkGap", len(txtSplit)/chunkGap)
		for i := 0; i <= len(txtSplit)/chunkGap; i++ {
			idxStart := i * chunkGap
			if idxStart > len(txtSplit) {
				continue
			}

			idxEnd := idxStart + chunkSize
			if idxEnd > len(txtSplit) {
				idxEnd = len(txtSplit)
			}
			// fmt.Println("---> idxStart, idxEnd, idxEnd-idxStart, chunkSize", idxStart, idxEnd, idxEnd-idxStart, chunkSize)

			// probably all parts already recorded on previous chunk
			if idxEnd-idxStart <= chunkOverlap {
				continue
			}

			chunkText := strings.Join(txtSplit[idxStart:idxEnd], " ")
			result = append(result, chunkText)
		}

		return result
	case config.ChunkModeParagraph:
		result := make([]string, 0)

		txtSplit := strings.Split(txt, "\n\n")
		chunkGap := chunkSize - chunkOverlap
		// fmt.Println("---> len(txtSplit)", len(txtSplit))
		// fmt.Println("---> chunkGap", chunkGap)
		// fmt.Println("---> len(txtSplit)/chunkGap", len(txtSplit)/chunkGap)
		for i := 0; i <= len(txtSplit)/chunkGap; i++ {
			idxStart := i * chunkGap
			if idxStart > len(txtSplit) {
				continue
			}

			idxEnd := idxStart + chunkSize
			if idxEnd > len(txtSplit) {
				idxEnd = len(txtSplit)
			}
			// fmt.Println("---> idxStart, idxEnd, idxEnd-idxStart, chunkSize", idxStart, idxEnd, idxEnd-idxStart, chunkSize)

			// probably all parts already recorded on previous chunk
			if idxEnd-idxStart <= chunkOverlap {
				continue
			}

			chunkText := strings.Join(txtSplit[idxStart:idxEnd], "\n\n")
			result = append(result, chunkText)
		}

		return result
	case config.ChunkModeParagraphWordSmart:
		result := make([]string, 0)

		txtSplit := strings.Split(txt, "\n\n")
		idxStart := 0

		for i := 0; i < len(txtSplit); i++ {
			chunkText := strings.Join(txtSplit[idxStart:i], "\n\n")

			if len(strings.Split(chunkText, " ")) < chunkSize {
				continue
			}

			idxStart = i
			for idxStart > 0 {
				if len(strings.Split(strings.Join(txtSplit[idxStart:i], "\n\n"), " ")) < chunkOverlap {
					idxStart--
					continue
				}
				break
			}

			// fmt.Println("---> idxStart, idxEnd", idxStart, i)
			result = append(result, chunkText)
		}

		chunkText := strings.Join(txtSplit[idxStart:], "\n\n")
		// fmt.Println("---> idxStart, idxEnd", idxStart, len(txtSplit))

		// probably all parts already recorded on previous chunk
		if len(strings.Split(chunkText, " ")) > chunkOverlap {
			result = append(result, chunkText)
		}

		return result
	default:
		return []string{txt}
	}

}
