package tdb

import "os"

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
