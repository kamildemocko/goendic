package main

import (
	"flag"
	"io"
	"log"
	"log/slog"
	"strings"

	"github.com/kamildemocko/goendic/internal/data"
	"github.com/kamildemocko/goendic/internal/logs"
	"github.com/kamildemocko/goendic/internal/printer"
	"github.com/kamildemocko/goendic/internal/repository"
	"github.com/kamildemocko/goendic/internal/repository/sqlite"
)

const (
	downloadUrl         = `https://en-word.net/static/english-wordnet-2024.xml.gz`
	searchedWordLogPath = `logs\words.log`
)

var (
	exactMatch bool
	allResults bool
	debugMode  bool
)

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

	printer.PrintFirstTimeDB()

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
	err := logs.InitLogger(searchedWordLogPath)
	if err != nil {
		panic(err)
	}

	flag.BoolVar(&exactMatch, "e", false, "use exact matching")
	flag.BoolVar(&allResults, "l", false, "return all results")
	flag.BoolVar(&debugMode, "d", false, "debug mode")
}

func main() {
	defer logs.CloseLogger()
	printer.SetupPrintUsage()
	flag.Parse()

	if !debugMode {
		log.SetOutput(io.Discard)
	}

	args := flag.Args()

	if len(args) < 1 {
		flag.Usage()
		return
	}

	searchedCompound := strings.Join(args, " ")
	slog.Info(searchedCompound)

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

	printer.PrintResult(results, allResults)
}
