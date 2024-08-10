package taskmanagement

import (
	"fmt"
	"log"
	"os"
)

// TODO: watch if database is deleted or modified

const (
	appFolderName = ".daily-term"
	databaseName  = "database.json"
)

func ensureAppFolderExists() {
	path := createPath()

	stat, err := os.Stat(path)

	if err == nil {
		return
	}

	if os.IsNotExist(err) {
		if err = os.Mkdir(path, 0777); err != nil {
			panic(err)
		}

		return
	}

	if !stat.IsDir() {
		panic(fmt.Sprintf(`"%v" is not a directory`, appFolderName))
	}
}

func CreateRepository() (*Repository, error) {
	ensureAppFolderExists()

	dbPath := createPath(databaseName)

	file, err := os.OpenFile(dbPath, os.O_CREATE|os.O_RDWR, 0600)

	if err != nil {
		return nil, err
	}

	return &Repository{
		file: file,
	}, nil
}

func (r *Repository) CloseRepository() {
	err := r.file.Close()

	if err != nil {
		log.Fatal(err)
	}
}
