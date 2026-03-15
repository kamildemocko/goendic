package main

import (
	"embed"
	"flag"
	"io"
	"log"
	"log/slog"
	"os"
	"strings"

	"github.com/kamildemocko/goendic/v2/internal/bootstrap"
	"github.com/kamildemocko/goendic/v2/internal/logs"
	"github.com/kamildemocko/goendic/v2/internal/printer"
	"github.com/kamildemocko/goendic/v2/internal/repository"
)

var (
	exactMatch  bool
	allResults  bool
	updateDb    bool
	debugMode   bool
	showVersion bool
	//go:embed version.txt
	versionFile embed.FS
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
	flag.BoolVar(&showVersion, "v", false, "show version")
}

func readVersion() string {
	data, err := versionFile.ReadFile("version.txt")
	if err != nil {
		return "dev"
	}
	return strings.TrimSpace(string(data))
}

func main() {
	defer logs.CloseLogger()
	printer.SetupPrintUsage()
	flag.Parse()

	if showVersion {
		printer.PrintVersion(readVersion())
		os.Exit(0)
	}

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
