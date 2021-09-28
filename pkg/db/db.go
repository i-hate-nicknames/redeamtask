package db

type PostgresConfig struct {
	DBName string
}

type BookRecord struct {
	Title, Author string
}

type BookDB interface {
	Save(*BookRecord) error
	Get(id int) (*BookRecord, error)
	Close() error
}

type postgresDB struct {
	// todo: postgres conn be here
}

func MakePostgresDB(conf PostgresConfig) BookDB {
	panic("not implemented")
}
