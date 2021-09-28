package api

import "github.com/i-hate-nicknames/redeamtask/pkg/db"

type BookAPI interface {
	Get() (*Book, error)
	Store(*Book) error
}

type Book db.BookRecord

func NewAPI(db db.BookDB) BookAPI {
	panic("not implemented")
}
