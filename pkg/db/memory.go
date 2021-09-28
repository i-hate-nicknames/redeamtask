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

func (db *memoryDB) Save(*BookRecord) error {
	panic("not implemented")
}

func (db *memoryDB) Get(id int) (*BookRecord, error) {
	panic("not implemented")
}

func (db *memoryDB) Close() error {
	return nil
}
