package main

import (
	"flag"
	"goendic/internal/data"
	"goendic/internal/printer"
	"goendic/internal/repository"
	"goendic/internal/repository/sqlite"
	"log"
	"strings"
)

const downloadUrl = `https://en-word.net/static/english-wordnet-2024.xml.gz`

var exactMatch bool

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

func init() {
	flag.BoolVar(&exactMatch, "e", false, "use exact matching")
}

func main() {
	flag.Parse()

	args := flag.Args()

	if len(args) < 1 {
		log.Println("Error: No search word provided")
		log.Println("Usage: endic [-e] WORD")
		log.Println(" -e  : Use exact matching")
		return
	}

	searchedCompound := strings.Join(args, " ")

	app := App{}
	repo, err := prepareData()
	if err != nil {
		panic(err)
	}
	app.repo = repo

	results, err := repo.FindWord(searchedCompound, exactMatch)
	if err != nil {
		panic(err)
	}

	if len(results) == 0 {
		printer.PrintEmpty()
		return
	}

	printer.PrintResult(results)
}
