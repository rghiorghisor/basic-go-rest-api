package util

import (
	"os"
	"path/filepath"
	"strings"
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

// SplitFilepath tries to split the provided path into the filename and directory.
// Note that no validation is done by this method.
func SplitFilepath(fileAndPath string) (string, string) {
	return filepath.Split(fileAndPath)
}

// ArrayToSetString converts the given string values to a set like struct.
func ArrayToSetString(values ...string) map[string]struct{} {
	set := make(map[string]struct{}, len(values))
	for _, s := range values {
		s = strings.TrimSpace(s)
		s = strings.ToLower(s)
		set[s] = struct{}{}
	}

	return set
}
