package api

import (
	"context"

	"github.com/i-hate-nicknames/redeamtask/pkg/book"
	"github.com/i-hate-nicknames/redeamtask/pkg/db"
)

// BookAPI implements domain level of book store
type BookAPI interface {
	Get(context.Context, int) (*book.Book, error)
	Create(context.Context, *book.Book) (int, error)
	Update(context.Context, int, *book.Book) error
	Delete(context.Context, int) error
}

type api struct {
	db db.BookDB
}

// NewAPI creates a new BookAPI backed by given database
func NewAPI(db db.BookDB) BookAPI {
	return &api{db: db}
}

func (api *api) Get(ctx context.Context, ID int) (*book.Book, error) {
	record, err := api.db.Get(ctx, ID)
	if err != nil {
		return nil, err
	}
	return recordToDomain(record), nil
}

func (api *api) Create(ctx context.Context, book *book.Book) (int, error) {
	dbBook, err := api.db.Create(ctx, domainToRecord(book))
	if err != nil {
		return 0, err
	}
	return dbBook.ID, nil
}

func (api *api) Update(ctx context.Context, bookID int, book *book.Book) error {
	record := domainToRecord(book)
	record.ID = bookID
	return api.db.Update(ctx, record)
}

func (api *api) Delete(ctx context.Context, bookID int) error {
	return api.db.Delete(ctx, bookID)
}

func recordToDomain(record *db.BookRecord) *book.Book {
	return &book.Book{
		Title:       record.Title,
		Publisher:   record.Publisher,
		Author:      record.Author,
		PublishDate: record.PublishDate,
		Status:      record.Status,
		Rating:      record.Rating,
	}
}

func domainToRecord(book *book.Book) *db.BookRecord {
	return &db.BookRecord{
		Book: book,
	}
}
