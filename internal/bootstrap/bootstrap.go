package bootstrap

import (
	"log"

	"github.com/kamildemocko/goendic/internal/data"
	"github.com/kamildemocko/goendic/internal/printer"
	"github.com/kamildemocko/goendic/internal/repository"
	"github.com/kamildemocko/goendic/internal/repository/sqlite"
)

func OpenRepo() (repository.Repository, error) {
	log.Println("opening repository")

	dsn, err := sqlite.CreateDBFileIfNotExists()
	if err != nil {
		return nil, err
	}

	repo, err := repository.InitSqliteDB(dsn)
	if err != nil {
		return nil, err
	}

	return repo, nil
}

func PrepareData(repo repository.Repository) error {
	log.Println("preparing data")

	dbHasData, err := repo.HasData()
	if err != nil {
		return err
	}

	if data.IsOffline() {
		log.Println("device is offline")
		if !dbHasData {
			printer.PrintOfflineState()
		}
		return nil
	}

	mostRecentUrl, err := data.FindMostRecentUrl()
	if err != nil {
		if !dbHasData {
		}
		return err
	}

	if !dbHasData {
		log.Println("first time update")
		printer.PrintFirstTimeDB()
		return updateDB(repo, mostRecentUrl)
	}

	currentDbUrl, err := repo.GetUrl()
	if err != nil {
		return err
	}

	if currentDbUrl != mostRecentUrl {
		printer.PrintOldDB()
	}

	return nil
}

func updateDB(repo repository.Repository, downloadUrl string) error {
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

	err = repo.UpdateUrl(downloadUrl)
	if err != nil {
		return err
	}

	return nil
}

func ForceUpdateDB(repo repository.Repository) error {
	log.Println("force updating database")

	if data.IsOffline() {
		log.Println("device is offline")
		printer.PrintOfflineState()
		return nil
	}

	printer.PrintUpdateDB()

	mostRecentUrl, err := data.FindMostRecentUrl()
	if err != nil {
		return err
	}

	err = updateDB(repo, mostRecentUrl)
	if err != nil {
		return err
	}

	printer.PrintDbUpdated()

	return nil
}
