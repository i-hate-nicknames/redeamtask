package db

import (
	"errors"

	"github.com/i-hate-nicknames/redeamtask/pkg/book"
)

type BookRecord struct {
	ID int // todo maybe uint
	*book.Book
}

type BookDB interface {
	Create(*BookRecord) (*BookRecord, error)
	Update(*BookRecord) error
	Get(id int) (*BookRecord, error)
	Delete(id int) error
	Close() error
	Migrate() error
}

var ErrBookNotFound = errors.New("book not found")
