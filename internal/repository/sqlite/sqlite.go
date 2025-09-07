package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"goendic/internal/data/model"
	"log"
	"os"
	"path/filepath"
	"time"
)

type SqliteRepository struct {
	DB *sql.DB
}

// creates DB file and returns DSN
func CreateDBFileIfNotExists() (string, error) {
	dbDir := "data"
	if err := os.MkdirAll(dbDir, 0755); err != nil {
		return "", err
	}

	path := filepath.Join(dbDir, "dictionary.db")
	dsn := fmt.Sprintf("file:%s?mode=rwc", path)

	return dsn, nil
}

func (sr *SqliteRepository) CreateTable() error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	log.Println("database init")

	query := `
	CREATE VIRTUAL TABLE IF NOT EXISTS dictionary USING fts5(
		word,
		pos,
		definition,
		examples
	);`

	_, err := sr.DB.ExecContext(ctx, query)
	if err != nil {
		return err
	}

	return nil
}

func (sr *SqliteRepository) UpdateData(entries []model.UpdateEntry) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	log.Println("updating database")

	tx, err := sr.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query_truncate := `
	DELETE FROM dictionary`

	_, err = tx.ExecContext(ctx, query_truncate)
	if err != nil {
		return err
	}

	query := `
	INSERT INTO dictionary (word, pos, definition, examples)
	VALUES (?, ?, ?, ?);`

	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, entry := range entries {
		_, err = stmt.ExecContext(
			ctx, entry.Word, entry.Pos, entry.Definition, entry.Examples,
		)
		if err != nil {
			return err
		}
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	log.Println("success")

	return nil
}

func (sr *SqliteRepository) FindWordExact(val string) error {
	return nil
}

func (sr *SqliteRepository) FindWord(val string) error {
	return nil
}
