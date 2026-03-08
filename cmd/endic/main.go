package main

import (
	"flag"
	"io"
	"log"
	"log/slog"
	"os"
	"strings"

	"github.com/kamildemocko/goendic/internal/bootstrap"
	"github.com/kamildemocko/goendic/internal/logs"
	"github.com/kamildemocko/goendic/internal/printer"
	"github.com/kamildemocko/goendic/internal/repository"
)

// TODO: dynamic URL
// TODO: check if new DB available, show a suggest with a new flag to update it

var (
	exactMatch bool
	allResults bool
	updateDb   bool
	debugMode  bool
)

type App struct {
	repo repository.Repository
}

func init() {
	err := logs.InitLogger()
	if err != nil {
		panic(err)
	}

	flag.BoolVar(&exactMatch, "e", false, "use exact matching")
	flag.BoolVar(&allResults, "l", false, "return all results")
	flag.BoolVar(&updateDb, "u", false, "update database")
	flag.BoolVar(&debugMode, "d", false, "debug mode")
}

func main() {
	defer logs.CloseLogger()
	printer.SetupPrintUsage()
	flag.Parse()

	if !debugMode {
		log.SetOutput(io.Discard)
	}

	repo, err := bootstrap.OpenRepo()
	if err != nil {
		panic(err)
	}

	if updateDb {
		bootstrap.ForceUpdateDB(repo)
		os.Exit(0)
	}

	word_args := flag.Args()
	if len(word_args) < 1 {
		flag.Usage()
		return
	}

	searchedCompound := strings.Join(word_args, " ")
	slog.Info(searchedCompound)

	app := App{}

	err = bootstrap.PrepareData(repo)
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
