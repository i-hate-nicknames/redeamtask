package api

import (
	"github.com/i-hate-nicknames/redeamtask/pkg/book"
	"github.com/i-hate-nicknames/redeamtask/pkg/db"
)

type BookAPI interface {
	Get() (*book.Book, error)
	Store(*book.Book) error
}

type api struct {
	db db.BookDB
}

func NewAPI(db db.BookDB) BookAPI {
	return &api{db: db}
}

func (api *api) Get() (*book.Book, error) {
	panic("not implemented")
}

func (api *api) Store(*book.Book) error {
	panic("not implemented")
}
