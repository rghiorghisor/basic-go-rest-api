package util

import (
	"os"
	"path/filepath"
)

// CreateParentFolder creates the parent folder for the provided file name. The
// folder will be created only if it does not exist.
func CreateParentFolder(fileName string) error {
	dirName := filepath.Dir(fileName)
	if _, statErr := os.Stat(dirName); statErr != nil {
		createErr := os.MkdirAll(dirName, os.ModePerm)
		if createErr != nil {
			return createErr
		}
	}

	return nil
}
