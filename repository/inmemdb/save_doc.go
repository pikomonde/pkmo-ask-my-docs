package inmemdb

import "pikomonde/ask-my-docs/entity"

func (r *InMemoryDBRepo) SaveDoc(key string, doc entity.Doc) error {
	r.db[key] = doc
	return nil
}
