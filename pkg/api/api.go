package api

import "github.com/i-hate-nicknames/redeamtask/pkg/db"

type BookAPI interface {
	Get() (*Book, error)
	Store(*Book) error
}

type Book db.BookRecord

type api struct {
	db db.BookDB
}

func NewAPI(db db.BookDB) BookAPI {
	return &api{db: db}
}

func (api *api) Get() (*Book, error) {
	panic("not implemented")
}

func (api *api) Store(*Book) error {
	panic("not implemented")
}
