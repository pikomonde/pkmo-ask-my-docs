package inmemdb

func (r *InMemoryDBRepo) SaveEmbedding(key string, doc []float32) error {
	r.db[key] = doc
	return nil
}
