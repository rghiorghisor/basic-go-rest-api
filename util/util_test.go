package util

import (
	"os"
	"testing"

	"gopkg.in/go-playground/assert.v1"
)

func TestCreateParentFolder(t *testing.T) {
	fileName := "../tests/testCreateDir/testFileName"
	dirName := "../tests/testCreateDir"

	CreateParentFolder(fileName)

	_, err := os.Stat(dirName)

	assert.Equal(t, false, os.IsNotExist(err))

	defer os.Remove(dirName)
}

func TestCreateParentFolderNoRights(t *testing.T) {
	fileName := "../tests/testCreateDir/tmp/testFileName"
	pName := "../tests/testCreateDir"
	dirName := "../tests/testCreateDir/tmp"

	CreateParentFolder(dirName)
	os.Chmod(pName, 0444)
	CreateParentFolder(fileName)
	_, err := os.Stat(dirName)

	assert.Equal(t, false, os.IsNotExist(err))

	defer func() {
		os.Chmod(pName, 0666)

		os.Remove(pName)
	}()
}

func TestArrayToSetString(t *testing.T) {
	assert.Equal(t, map[string]struct{}{}, ArrayToSetString())
	assert.Equal(t, map[string]struct{}{"a": struct{}{}, "b": struct{}{}}, ArrayToSetString("A ", "  b  "))
	assert.Equal(t, map[string]struct{}{"a": struct{}{}, "b": struct{}{}, "c": struct{}{}}, ArrayToSetString("A ", "  b  ", "C"))
}
