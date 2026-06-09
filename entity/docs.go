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
