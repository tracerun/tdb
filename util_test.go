package tdb

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFolderActions(t *testing.T) {
	demoFolder := "util_test_folder"
	demoFile := "util_test_file"

	exist, err := checkFolderExist(demoFolder)
	assert.NoError(t, err, "should have no error")
	assert.False(t, exist, "folder should not exist")

	var f *os.File
	f, err = os.Create(demoFile)
	if assert.NoError(t, err, "should have no error to create file") {
		defer func() {
			f.Close()
			err := os.Remove(demoFile)
			assert.NoError(t, err, "should have no error to delete file")
		}()
	}

	exist, err = checkFolderExist(demoFile)
	assert.Equal(t, ErrPathNotFolder, err, "should have error when path is not folder")
	assert.False(t, exist, "should return false if it is a file")

	err = createFolder(demoFolder)
	if assert.NoError(t, err, "should have no error to create folder") {
		defer func() {
			err := os.RemoveAll(demoFolder)
			assert.NoError(t, err, "should have no error to delete folder")
		}()
	}

	err = createFolder(demoFolder)
	assert.NoError(t, err, "should have no error if folder exist")

	exist, err = checkFolderExist(demoFolder)
	assert.NoError(t, err, "should have no error when path is a folder")
	assert.True(t, exist, "should return true if folder exist")
}

func TestListFolders(t *testing.T) {

}
