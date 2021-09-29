package db

import (
	"context"

	"github.com/jackc/pgx/v4"
)

type postgresDB struct {
	conn *pgx.Conn
}

func MakePostgresDB(ctx context.Context, dsn string) (BookDB, error) {
	conn, err := pgx.Connect(ctx, dsn)
	if err != nil {
		return nil, err
	}
	return &postgresDB{conn: conn}, nil
}

func (db *postgresDB) Migrate() error {
	panic("not implemented")
}

func (db *postgresDB) Create(br *BookRecord) (*BookRecord, error) {
	panic("not implemented")
}

func (db *postgresDB) Update(br *BookRecord) error {
	panic("not implemented")
}

func (db *postgresDB) Get(id int) (*BookRecord, error) {
	panic("not implemented")
}

func (db *postgresDB) Delete(ID int) error {
	panic("not implemented")
}

func (db *postgresDB) Close() error {
	panic("not implemented")
}
