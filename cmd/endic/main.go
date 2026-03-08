package main

import (
	"flag"
	"io"
	"log"
	"log/slog"
	"strings"

	"github.com/kamildemocko/goendic/internal/bootstrap"
	"github.com/kamildemocko/goendic/internal/logs"
	"github.com/kamildemocko/goendic/internal/printer"
	"github.com/kamildemocko/goendic/internal/repository"
)

const downloadUrl = `https://en-word.net/static/english-wordnet-2024.xml.gz`

var (
	exactMatch bool
	allResults bool
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
	repo, err := bootstrap.PrepareData(downloadUrl)
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
