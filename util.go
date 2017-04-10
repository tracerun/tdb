package tdb

import (
	"io/ioutil"
	"os"
)

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

// list all the folders in a given folder
func listFolders(folder string) ([]string, error) {
	files, err := ioutil.ReadDir(folder)
	if err != nil {
		return nil, err
	}

	var folders []string
	for _, f := range files {
		if f.IsDir() {
			folders = append(folders, f.Name())
		}
	}
	return folders, nil
}
