package taskmanagement

import (
	"log"
	"os"
)

// TODO: watch if database is deleted or modified

const databaseName = "database.json"

func CreateRepository() (*Repository, error) {
	file, err := os.OpenFile(databaseName, os.O_CREATE|os.O_RDWR, 0600)

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
