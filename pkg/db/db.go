package db

type PostgresConfig struct {
	DBName string
}

type BookRecord struct {
	id            int // todo maybe uint
	Title, Author string
}

type BookDB interface {
	Save(*BookRecord) error
	Get(id int) (*BookRecord, error)
	Close() error
}
