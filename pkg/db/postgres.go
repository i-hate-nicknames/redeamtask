package db

import (
	"context"
	"errors"
	"io/ioutil"
	"log"
	"os"

	"github.com/jackc/pgx/v4"
)

type postgresDB struct {
	conn *pgx.Conn
}

// MakePostgresDB creates a new postgres database out of given DSN
func MakePostgresDB(ctx context.Context, dsn string) (BookDB, error) {
	conn, err := pgx.Connect(ctx, dsn)
	if err != nil {
		return nil, err
	}
	return &postgresDB{conn: conn}, nil
}

func (db *postgresDB) Migrate(ctx context.Context) error {
	// todo: a migration library would be nice
	f, err := os.Open("sql/1-books.sql")
	if err != nil {
		return err
	}
	migration, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}
	_, err = db.conn.Exec(context.TODO(), string(migration))
	return err
}

func (db *postgresDB) Create(ctx context.Context, br *BookRecord) (*BookRecord, error) {
	query := `
		INSERT INTO books (title, author, publisher, publish_date, rating, _status)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`
	lastInsertedID := 0
	err := db.conn.QueryRow(ctx, query, br.Title, br.Author, br.Publisher, br.PublishDate, br.Rating, br.Status).Scan(&lastInsertedID)
	if err != nil {
		return nil, err
	}
	br.ID = lastInsertedID
	return br, nil
}

func (db *postgresDB) Update(ctx context.Context, br *BookRecord) error {
	query := `
		UPDATE books
		SET
			title = $1,
			author = $2,
			publisher = $3,
			publish_date = $4,
			rating = $5,
			_status = $6
		WHERE id = $7 AND deleted_at IS NULL
	`
	cmd, err := db.conn.Exec(ctx, query, br.Title, br.Author, br.Publisher, br.PublishDate, br.Rating, br.Status, br.ID)
	affected := cmd.RowsAffected()
	if affected == 0 {
		return ErrBookNotFound
	}
	if affected > 1 {
		log.Println("Single book update affected more than a single book")
	}
	return err
}

func (db *postgresDB) Get(ctx context.Context, ID int) (*BookRecord, error) {
	query := `
		SELECT id, title, author, publisher, publish_date, rating, _status
		FROM books
		WHERE id = $1 AND deleted_at IS NULL
	`
	var br BookRecord
	row := db.conn.QueryRow(ctx, query, ID)
	err := row.Scan(&br.ID, &br.Title, &br.Author, &br.Publisher, &br.PublishDate, &br.Rating, &br.Status)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrBookNotFound
	}
	if err != nil {
		return nil, err
	}
	return &br, nil
}

func (db *postgresDB) Delete(ctx context.Context, ID int) error {
	query := `
		UPDATE books
		SET deleted_at = now()
		WHERE id = $1 AND deleted_at IS NULL
	`
	cmd, err := db.conn.Exec(ctx, query, ID)
	affected := cmd.RowsAffected()
	if affected > 1 {
		log.Println("Single book delete affected more than a single book")
	}
	return err
}

func (db *postgresDB) GetAll(_ context.Context) ([]*BookRecord, error) {
	panic("not implemented")
}

func (db *postgresDB) Close(ctx context.Context) error {
	return db.conn.Close(ctx)
}
