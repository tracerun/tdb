package tdb

import (
	"errors"
)

var (
	// ErrDBPathNotFolder path for TDB is not a folder
	ErrDBPathNotFolder = errors.New("path for TDB is not a folder")
	// ErrInfoFilePath path for info file is not a file
	ErrInfoFilePath = errors.New("path for info file is not a file")
	// ErrNotExist something not exist
	ErrNotExist = errors.New("not exist")
)
