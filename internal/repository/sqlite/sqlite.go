package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/kamildemocko/goendic/v2/internal/data/model"
)

type SqliteRepository struct {
	DB *sql.DB
}

type searchResult struct {
	entry    model.UpdateEntry
	distance int
	isExact  bool
}

func levenshteinDistance(s1, s2 string) int {
	s1Lower := strings.ToLower(s1)
	s2Lower := strings.ToLower(s2)

	if len(s1Lower) == 0 {
		return len(s2Lower)
	}
	if len(s2Lower) == 0 {
		return len(s1Lower)
	}

	// Create matrix
	matrix := make([][]int, len(s1Lower)+1)
	for i := range matrix {
		matrix[i] = make([]int, len(s2Lower)+1)
	}

	// Initialize first column and row
	for i := 0; i <= len(s1Lower); i++ {
		matrix[i][0] = i
	}
	for j := 0; j <= len(s2Lower); j++ {
		matrix[0][j] = j
	}

	// Fill matrix
	for i := 1; i <= len(s1Lower); i++ {
		for j := 1; j <= len(s2Lower); j++ {
			cost := 0
			if s1Lower[i-1] != s2Lower[j-1] {
				cost = 1
			}

			matrix[i][j] = min(
				matrix[i-1][j]+1,      // deletion
				matrix[i][j-1]+1,      // insertion
				matrix[i-1][j-1]+cost, // substitution
			)
		}
	}

	return matrix[len(s1Lower)][len(s2Lower)]
}

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
	CREATE TABLE IF NOT EXISTS dictionary (
		word TEXT NOT NULL,
		pos TEXT,
		definition TEXT,
		examples TEXT
	);
	CREATE INDEX IF NOT EXISTS idx_dictionary_word ON dictionary(word);
	CREATE INDEX IF NOT EXISTS idx_dictionary_word_lower ON dictionary(lower(word));`

	_, err := sr.DB.ExecContext(ctx, query)
	if err != nil {
		return err
	}

	query_url := `
	CREATE TABLE IF NOT EXISTS url (
		value TEXT NOT NULL
	);`

	_, err = sr.DB.ExecContext(ctx, query_url)
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

func (sr *SqliteRepository) UpdateUrl(url string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	log.Println("updating url")

	tx, err := sr.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query_truncate := `
	DELETE FROM url`

	_, err = tx.ExecContext(ctx, query_truncate)
	if err != nil {
		return err
	}

	query := `
	INSERT INTO url (value)
	VALUES (?);`

	_, err = tx.ExecContext(ctx, query, url)
	if err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	log.Println("success")
	return nil
}

func (sr *SqliteRepository) GetUrl() (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	query := `SELECT value FROM url`

	var value string
	err := sr.DB.QueryRowContext(ctx, query).Scan(&value)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", nil
		}
		return "", err
	}

	return value, nil
}

func (sr *SqliteRepository) fuzzySearch(ctx context.Context, val string, seen map[string]bool) ([]model.UpdateEntry, error) {
	query := `SELECT word, pos, definition, examples FROM dictionary`

	rows, err := sr.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var candidates []searchResult
	maxDistance := 3

	for rows.Next() {
		var entry model.UpdateEntry
		err := rows.Scan(&entry.Word, &entry.Pos, &entry.Definition, &entry.Examples)
		if err != nil {
			return nil, err
		}

		if seen[strings.ToLower(entry.Word)] {
			continue
		}

		firstWord := entry.Word
		if idx := strings.Index(entry.Word, " "); idx > 0 {
			firstWord = entry.Word[:idx]
		}

		distance := levenshteinDistance(val, firstWord)
		if distance <= maxDistance {
			candidates = append(candidates, searchResult{
				entry:    entry,
				distance: distance,
				isExact:  distance == 0,
			})
		}
	}

	// Sort by distance, then by length, then alphabetically
	sort.Slice(candidates, func(i, j int) bool {
		if candidates[i].distance != candidates[j].distance {
			return candidates[i].distance < candidates[j].distance
		}
		if len(candidates[i].entry.Word) != len(candidates[j].entry.Word) {
			return len(candidates[i].entry.Word) < len(candidates[j].entry.Word)
		}
		return candidates[i].entry.Word < candidates[j].entry.Word
	})

	maxResults := 100
	if len(candidates) > maxResults {
		candidates = candidates[:maxResults]
	}

	var results []model.UpdateEntry
	for _, c := range candidates {
		results = append(results, c.entry)
	}

	return results, nil
}

func (sr *SqliteRepository) FindWord(val string, exact bool) ([]model.UpdateEntry, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if exact {
		query := `
		SELECT word, pos, definition, examples
        FROM dictionary
        WHERE lower(word) = lower(?)`

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

		if len(results) > 0 {
			return results, nil
		}

		// Fall back to fuzzy search if exact match found nothing
		return sr.fuzzySearch(ctx, val, make(map[string]bool))

	} else {
		// First try prefix matching and similar words
		query := `
		SELECT word, pos, definition, examples
        FROM dictionary
        WHERE lower(word) = lower(?)
           OR lower(word) = lower(?) || 's'
           OR lower(word) LIKE lower(?) || ' %'
           OR lower(word) LIKE lower(?) || '%'
        ORDER BY
           CASE WHEN lower(word) = lower(?) THEN 1
                WHEN lower(word) = lower(?) || 's' THEN 2
                WHEN lower(word) LIKE lower(?) || ' %' THEN 3
                ELSE 4
           END,
           length(word) ASC,
           word ASC
        LIMIT 100`

		rows, err := sr.DB.QueryContext(ctx, query, val, val, val, val, val, val, val)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		var results []model.UpdateEntry
		seen := make(map[string]bool)
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
			seen[strings.ToLower(entry.Word)] = true
		}

		if len(results) > 0 {
			return results, nil
		}

		// Fall back to fuzzy search with Levenshtein distance
		return sr.fuzzySearch(ctx, val, seen)
	}
}
