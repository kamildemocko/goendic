package repository

import (
	"database/sql"
	"endic/internal/data/model"
	"endic/internal/repository/sqlite"
	"log"
	"time"

	_ "modernc.org/sqlite"
)

type Repository interface {
	CreateTable() error
	HasData() (bool, error)
	UpdateData([]model.UpdateEntry) error
	FindWord(val string, exact bool) ([]model.UpdateEntry, error)
}

func NewSqliteDB(db *sql.DB) Repository {
	return &sqlite.SqliteRepository{
		DB: db,
	}
}

func InitSqliteDB(dsn string) (Repository, error) {
	log.Println("connecting to DB")

	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	db.SetConnMaxLifetime(5 * time.Minute)

	repo := NewSqliteDB(db)
	err = repo.CreateTable()
	if err != nil {
		return nil, err
	}

	return repo, nil
}
