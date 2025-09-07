package main

import (
	"goendic/internal/data"
	"goendic/internal/repository"
	"goendic/internal/repository/sqlite"
	"log"
)

const downloadUrl = `https://en-word.net/static/english-wordnet-2024.xml.gz`

type App struct {
	repo repository.Repository
}

func prepareData() (repository.Repository, error) {
	dsn, err := sqlite.CreateDBFileIfNotExists()
	if err != nil {
		return nil, err
	}

	repo, err := repository.InitSqliteDB(dsn)
	if err != nil {
		return nil, err
	}

	dbExists, err := repo.HasData()
	if err != nil {
		return nil, err
	}
	if dbExists {
		return repo, nil
	}

	loader := data.NewDataLoader(downloadUrl)
	file, err := loader.Get()
	if err != nil {
		return nil, err
	}
	defer loader.Close()

	data, err := data.ParseXML(file)
	if err != nil {
		return nil, err
	}

	err = repo.UpdateData(data)
	if err != nil {
		return nil, err
	}

	return repo, nil
}

func main() {
	app := App{}
	repo, err := prepareData()
	if err != nil {
		panic(err)
	}
	app.repo = repo

	searchedWord := "ambidex"

	// job
	results, err := repo.FindWord(searchedWord)
	if err != nil {
		panic(err)
	}
	for _, result := range results {
		log.Printf("%+v", result)
	}
}
