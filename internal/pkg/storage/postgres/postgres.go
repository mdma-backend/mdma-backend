package postgres

import (
	"database/sql"
	"fmt"
	"log"

	_ "embed"

	_ "github.com/lib/pq"
)

//go:embed schema.sql
var schema string

type DB struct {
	pool *sql.DB
}

func New(dsn string) (*DB, error) {
	pool, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("sql open: %w", err)
	}

	if err := pool.Ping(); err != nil {
		return nil, fmt.Errorf("ping: %w", err)
	}

	if err := migrateDatabase(pool); err != nil {
		return nil, fmt.Errorf("migrate db: %w", err)
	}

	return &DB{
		pool: pool,
	}, nil
}

func migrateDatabase(pool *sql.DB) error {
	var tableCount uint32
	if err := pool.QueryRow(`
SELECT COUNT(table_name)
FROM information_schema.tables
WHERE table_schema = 'public'
`).Scan(&tableCount); err != nil {
		return err
	}

	if tableCount > 0 {
		log.Printf("database already populated with %d tables", tableCount)
		return nil
	}

	log.Println("initializing new database")

	if _, err := pool.Exec(schema); err != nil {
		return err
	}

	return nil
}
