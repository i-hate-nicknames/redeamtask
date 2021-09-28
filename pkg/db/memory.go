package db

import "sync"

type memoryDB struct {
	nextID int
	items  map[int]*BookRecord
	mu     sync.RWMutex
}

func MakeMemoryDB() BookDB {
	items := make(map[int]*BookRecord)
	db := &memoryDB{items: items, nextID: 1}
	return db
}

func (db *memoryDB) Save(br *BookRecord) error {
	db.mu.Lock()
	defer db.mu.Unlock()
	db.items[br.id] = br
	return nil
}

func (db *memoryDB) Get(id int) (*BookRecord, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()
	b, ok := db.items[id]
	if !ok {
		return nil, ErrBookNotFound
	}
	return b, nil
}

func (db *memoryDB) Close() error {
	return nil
}
