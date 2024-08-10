package taskmanagement

import (
	"os"
	"path"
)

func createPath(chunks ...string) string {
	homeDir, err := os.UserHomeDir()

	if err != nil {
		panic(err)
	}

	fullPath := []string{homeDir, appFolderName}
	fullPath = append(fullPath, chunks...)

	return path.Join(fullPath...)
}
