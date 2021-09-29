package db

import (
	"context"
	"errors"

	"github.com/i-hate-nicknames/redeamtask/pkg/book"
)

// BookRecord represents a single record in the book database
type BookRecord struct {
	ID int // todo maybe uint
	*book.Book
}

// BookDB is a generic book database interface
type BookDB interface {
	Create(context.Context, *BookRecord) (*BookRecord, error)
	Update(context.Context, *BookRecord) error
	Get(context.Context, int) (*BookRecord, error)
	GetAll(context.Context) ([]*BookRecord, error)
	Delete(context.Context, int) error
	Close(context.Context) error
	Migrate(context.Context) error
}

// ErrBookNotFound is returned when requested book is not present in the database
var ErrBookNotFound = errors.New("book not found")
