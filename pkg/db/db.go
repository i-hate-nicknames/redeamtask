package db

import (
	"context"
	"errors"

	"github.com/i-hate-nicknames/redeamtask/pkg/book"
)

type BookRecord struct {
	ID int // todo maybe uint
	*book.Book
}

type BookDB interface {
	Create(context.Context, *BookRecord) (*BookRecord, error)
	Update(context.Context, *BookRecord) error
	Get(context.Context, int) (*BookRecord, error)
	Delete(context.Context, int) error
	Close(context.Context) error
	Migrate(context.Context) error
}

var ErrBookNotFound = errors.New("book not found")
