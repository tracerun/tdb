package tdb

import (
	"math/rand"
	"os"
	"path/filepath"
	"strings"
)

const (
	letterBytes = "0123456789abcdefghijklmnopqrstuvwxyz"
	letterLen   = 36
)

func randBytes(n int) []byte {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(letterLen)]
	}
	return b
}

// checkFolderExist
// return an error if folder is a file.
func checkFolderExist(folder string) (bool, error) {
	if stat, err := os.Stat(folder); err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	} else if !stat.IsDir() {
		return false, ErrPathNotFolder
	}
	return true, nil
}

// createFolder to create a folder, if exist, do nothing
// return an error if folder is a file.
func createFolder(folder string) error {
	exist, err := checkFolderExist(folder)
	if err != nil {
		return err
	}

	if !exist {
		err = os.MkdirAll(folder, os.ModePerm)
	}
	return err
}

// get the file name without ext and path
func getFileName(fPath string) string {
	baseName := filepath.Base(fPath)
	return strings.TrimSuffix(baseName, filepath.Ext(baseName))
}
