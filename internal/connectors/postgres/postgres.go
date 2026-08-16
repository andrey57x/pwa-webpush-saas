package postgres

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib" // Регистрирует PGX драйвер для стандартного database/sql
)

func NewPool(dsn string) (*sql.DB, error) {
	// Инициализируем стандартный *sql.DB через драйвер pgx
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database/sql connection: %w", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(30 * time.Minute)
	db.SetConnMaxIdleTime(5 * time.Minute)

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database/sql: %w", err)
	}

	log.Println("Successfully connected to Database via standard database/sql!")
	return db, nil
}
