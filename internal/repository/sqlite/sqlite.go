package sqlite

import (
	"context"
	"database/sql"
	"time"
)

type SqliteRepository struct {
	DB *sql.DB
}

func (sr *SqliteRepository) CreateTable() error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

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

func (sr *SqliteRepository) UpdateData() error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	query := `
	INSERT INTO dictionary (word, pos, definition, examples)
	VALUES (
		'bomb',
		'noun',
		'an explosive device fused to explode under specific conditions',
		'The army diffused a bomb.'
	);`

	_, err := sr.DB.ExecContext(ctx, query)
	if err != nil {
		return err
	}

	return nil
}

func (sr *SqliteRepository) FindWordExact(val string) error {
	return nil
}

func (sr *SqliteRepository) FindWord(val string) error {
	return nil
}
