package sqlite

import (
	"context"
	"database/sql"
	"endic/internal/data/model"
	"fmt"
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
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	dbDir := filepath.Join(configDir, "goendic")

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

func (sr *SqliteRepository) HasData() (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	queryTable := `
	SELECT name FROM sqlite_master
	WHERE type='table' AND name='dictionary'`

	var tableName string
	err := sr.DB.QueryRowContext(ctx, queryTable).Scan(&tableName)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	queryCount := `
	SELECT COUNT(*) FROM dictionary LIMIT 1`

	var count int
	err = sr.DB.QueryRowContext(ctx, queryCount).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
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

func (sr *SqliteRepository) FindWord(val string, exact bool) ([]model.UpdateEntry, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	query := `
	SELECT word, pos, definition, examples 
	FROM dictionary
	WHERE word MATCH ?`

	if !exact {
		val = val + "*"
	}

	rows, err := sr.DB.QueryContext(ctx, query, val)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []model.UpdateEntry

	for rows.Next() {
		var entry model.UpdateEntry

		err := rows.Scan(
			&entry.Word,
			&entry.Pos,
			&entry.Definition,
			&entry.Examples,
		)
		if err != nil {
			return nil, err
		}

		results = append(results, entry)
	}

	return results, nil
}
