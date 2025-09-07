package main

import (
	"goendic/internal/data"
	"goendic/internal/repository"
	"goendic/internal/repository/sqlite"
)

const downloadUrl = `https://en-word.net/static/english-wordnet-2024.xml.gz`

func prepareData() error {
	dsn, err := sqlite.CreateDBFileIfNotExists()
	if err != nil {
		return err
	}

	repo, err := repository.InitSqliteDB(dsn)
	if err != nil {
		return err
	}

	dbExists, err := repo.HasData()
	if err != nil {
		return err
	}
	if dbExists {
		return nil
	}

	loader := data.NewDataLoader(downloadUrl)
	file, err := loader.Get()
	if err != nil {
		return err
	}
	defer loader.Close()

	data, err := data.ParseXML(file)
	if err != nil {
		return err
	}

	err = repo.UpdateData(data)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	err := prepareData()
	if err != nil {
		panic(err)
	}

	// job
}
