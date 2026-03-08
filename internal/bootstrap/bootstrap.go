package bootstrap

import (
	"github.com/kamildemocko/goendic/internal/data"
	"github.com/kamildemocko/goendic/internal/printer"
	"github.com/kamildemocko/goendic/internal/repository"
	"github.com/kamildemocko/goendic/internal/repository/sqlite"
)

func PrepareData(dbUrl string) (repository.Repository, error) {
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

	loader := data.NewDataLoader(dbUrl)
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
