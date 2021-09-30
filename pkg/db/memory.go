package db

import (
	"context"
	"sync"

	"github.com/rs/zerolog"
)

type memoryDB struct {
	nextID int
	items  map[int]BookRecord
	mu     sync.RWMutex
	logger zerolog.Logger
}

// MakeMemoryDB creates a new in-memory database
func MakeMemoryDB(logger zerolog.Logger) BookDB {
	items := make(map[int]BookRecord)
	db := &memoryDB{items: items, nextID: 1, logger: logger}
	return db
}

func (db *memoryDB) save(br BookRecord) error {
	db.mu.Lock()
	defer db.mu.Unlock()
	db.items[br.ID] = br
	return nil
}

func (db *memoryDB) Migrate(_ context.Context) error {
	return nil
}

func (db *memoryDB) Create(_ context.Context, br BookRecord) (BookRecord, error) {
	db.logger.Debug().Str("book_title", br.Title).Msg("Create book")
	db.mu.Lock()
	br.ID = db.nextID
	db.nextID++
	db.mu.Unlock()
	err := db.save(br)
	if err != nil {
		return BookRecord{}, err
	}
	return br, nil
}

func (db *memoryDB) Update(ctx context.Context, br BookRecord) error {
	db.logger.Debug().Int("book_id", br.ID).Msg("Update book")
	if _, err := db.Get(ctx, br.ID); err != nil {
		return err
	}
	return db.save(br)
}

func (db *memoryDB) Get(_ context.Context, ID int) (BookRecord, error) {
	db.logger.Debug().Int("book_id", ID).Msg("Get book")
	db.mu.RLock()
	defer db.mu.RUnlock()
	b, ok := db.items[ID]
	if !ok {
		return BookRecord{}, ErrBookNotFound
	}
	return b, nil
}

func (db *memoryDB) GetAll(_ context.Context) ([]BookRecord, error) {
	db.logger.Debug().Msg("Get all books")
	db.mu.RLock()
	defer db.mu.RUnlock()
	var records []BookRecord
	for _, item := range db.items {
		records = append(records, item)
	}
	return records, nil
}

func (db *memoryDB) Delete(_ context.Context, ID int) error {
	db.logger.Debug().Int("book_id", ID).Msg("Delete book")
	db.mu.Lock()
	defer db.mu.Unlock()
	delete(db.items, ID)
	return nil
}

func (db *memoryDB) Close(_ context.Context) error {
	return nil
}
