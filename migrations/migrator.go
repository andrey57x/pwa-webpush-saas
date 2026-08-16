package migrations

import (
	"embed"
	"errors"
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

// Вшиваем все .sql файлы из папки migrations в бинарник Go
//
//go:embed *.sql
var migrationFiles embed.FS

// RunMigrations автоматически проверяет и применяет новые SQL-миграции в БД
func RunMigrations(databaseURL string) error {
	log.Println("Checking database migrations...")

	// 1. Инициализируем iofs источник из embedded файлов
	sourceDriver, err := iofs.New(migrationFiles, ".")
	if err != nil {
		return fmt.Errorf("failed to create iofs migration source: %w", err)
	}

	// 2. Создаем инстанс мигратора
	m, err := migrate.NewWithSourceInstance("iofs", sourceDriver, databaseURL)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}
	defer m.Close()

	// 3. Применяем миграции UP
	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			log.Println("Database is up to date (no new migrations to apply).")
			return nil
		}
		return fmt.Errorf("failed to apply migrations: %w", err)
	}

	log.Println("Successfully applied database migrations!")
	return nil
}
