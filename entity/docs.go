package entity

import (
	"fmt"
)

type Docs []Doc

type Doc struct {
	Path             string
	ChunkNo          int
	Text             string
	Embeddings       []float32
	CosineSimilarity float32
}

func (d Doc) Key() string {
	return fmt.Sprintf("%s::%d", d.Path, d.ChunkNo)
}

// func MustParseDocKey(key string) (path string, chunkNo int) {
// 	parts := strings.Split(key, "::")
// 	if len(parts) != 2 {
// 		return "", 0
// 	}

// 	chunkNo, err := strconv.Atoi(parts[1])
// 	if err != nil {
// 		return "", 0
// 	}

// 	// d.Path, d.ChunkNo = parts[0], chunkNo

// 	return parts[0], chunkNo
// }
