package main

import (
	"goendic/internal/data"
	"goendic/internal/repository"
	"goendic/internal/repository/sqlite"
)

const downloadUrl = `https://en-word.net/static/english-wordnet-2024.xml.gz`

func main() {
	dsn, err := sqlite.CreateDBFileIfNotExists()
	if err != nil {
		panic(err)
	}

	repo, err := repository.InitSqliteDB(dsn)
	if err != nil {
		panic(err)
	}

	loader := data.NewDataLoader(downloadUrl)
	file, err := loader.Get()
	if err != nil {
		panic(err)
	}
	defer loader.Close()

	data, err := data.ParseXML(file)
	if err != nil {
		panic(err)
	}

	_ = repo.UpdateData(data)
}
