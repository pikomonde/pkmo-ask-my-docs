package inmemdb

import (
	"errors"
	"fmt"
	"math"
	"pikomonde/ask-my-docs/entity"
	"sort"
)

func (r *InMemoryDBRepo) Search(queryVector []float32) ([]entity.Doc, error) {
	res := make([]entity.Doc, 0)

	// calculate cosine similarity
	for key, doc := range r.db {
		similiarity, err := cosineSimilarity(doc.Embeddings, queryVector)
		if err != nil {
			return nil, fmt.Errorf("error when calculate cosine similarity")
		}

		doc.CosineSimilarity = similiarity
		res = append(res, doc)

		// save to db
		r.db[key] = doc
	}

	// sort by similarity
	sort.Slice(res, func(i, j int) bool {
		return res[i].CosineSimilarity > res[j].CosineSimilarity
	})

	return res, nil
}

func cosineSimilarity(a, b []float32) (float32, error) {
	if len(a) != len(b) {
		return 0, errors.New("vectors must have the same dimensions")
	}
	if len(a) == 0 {
		return 0, errors.New("vectors cannot be empty")
	}

	var dotProduct, normA, normB float32

	for i := range a {
		dotProduct += a[i] * b[i]
		normA += a[i] * a[i]
		normB += b[i] * b[i]
	}

	// prevent divide by zero error
	if normA == 0 || normB == 0 {
		return 0, nil
	}

	return dotProduct / (float32(math.Sqrt(float64(normA))) * float32(math.Sqrt(float64(normB)))), nil
}
