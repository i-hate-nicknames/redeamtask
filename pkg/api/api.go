package api

import (
	"github.com/i-hate-nicknames/redeamtask/pkg/book"
	"github.com/i-hate-nicknames/redeamtask/pkg/db"
)

type BookAPI interface {
	Get(id int) (*book.Book, error)
	Create(*book.Book) (int, error)
	Update(int, *book.Book) error
}

type api struct {
	db db.BookDB
}

func NewAPI(db db.BookDB) BookAPI {
	return &api{db: db}
}

func (api *api) Get(id int) (*book.Book, error) {
	record, err := api.db.Get(id)
	if err != nil {
		return nil, err
	}
	return recordToDomain(record), nil
}

func (api *api) Create(book *book.Book) (int, error) {
	dbBook, err := api.db.Create(domainToRecord(book))
	if err != nil {
		return 0, err
	}
	return dbBook.ID, nil
}

func (api *api) Update(bookID int, book *book.Book) error {
	record := domainToRecord(book)
	record.ID = bookID
	return api.db.Update(record)
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
