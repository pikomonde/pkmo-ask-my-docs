package inmemdb

import "pikomonde/ask-my-docs/entity"

func (r *InMemoryDBRepo) Search(queryVector []float32) (entity.Doc, error) {
	return entity.Doc{}, nil
}
