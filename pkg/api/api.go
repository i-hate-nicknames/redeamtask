package api

import (
	"github.com/i-hate-nicknames/redeamtask/pkg/book"
	"github.com/i-hate-nicknames/redeamtask/pkg/db"
)

type BookAPI interface {
	Get(id int) (*book.Book, error)
	Create(*book.Book) error
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

// todo: think about design of create/update operations

func (api *api) Create(book *book.Book) error {
	return api.db.Create(domainToRecord(book))
}

func (api *api) Update(bookID int, book *book.Book) error {
	record := domainToRecord(book)
	record.ID = bookID
	return api.db.Update(record)
}

func recordToDomain(record *db.BookRecord) *book.Book {
	panic("not implemented")
}

func domainToRecord(book *book.Book) *db.BookRecord {
	panic("not implemented")
}
