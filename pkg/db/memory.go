package db

import (
	"context"
	"sync"
)

type memoryDB struct {
	nextID int
	items  map[int]*BookRecord
	mu     sync.RWMutex
}

// MakeMemoryDB creates a new in-memory database
func MakeMemoryDB() BookDB {
	items := make(map[int]*BookRecord)
	db := &memoryDB{items: items, nextID: 1}
	return db
}

func (db *memoryDB) save(br *BookRecord) error {
	db.mu.Lock()
	defer db.mu.Unlock()
	db.items[br.ID] = br
	return nil
}

func (db *memoryDB) Migrate(_ context.Context) error {
	return nil
}

func (db *memoryDB) Create(_ context.Context, br *BookRecord) (*BookRecord, error) {
	db.mu.Lock()
	br.ID = db.nextID
	db.nextID++
	db.mu.Unlock()
	err := db.save(br)
	if err != nil {
		return nil, err
	}
	return br, nil
}

func (db *memoryDB) Update(ctx context.Context, br *BookRecord) error {
	if _, err := db.Get(ctx, br.ID); err != nil {
		return err
	}
	return db.save(br)
}

func (db *memoryDB) Get(_ context.Context, ID int) (*BookRecord, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()
	b, ok := db.items[ID]
	if !ok {
		return nil, ErrBookNotFound
	}
	return b, nil
}

func (db *memoryDB) Delete(_ context.Context, ID int) error {
	db.mu.Lock()
	defer db.mu.Unlock()
	delete(db.items, ID)
	return nil
}

func (db *memoryDB) Close(_ context.Context) error {
	return nil
}
