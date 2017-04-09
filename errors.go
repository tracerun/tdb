package tdb

import (
	"errors"
)

var (
	// ErrPathNotFolder path is not a folder
	ErrPathNotFolder = errors.New("path for is not a folder")
	// ErrPathNotFile path is not a file
	ErrPathNotFile = errors.New("path is not a file")
	// ErrNotExist something not exist
	ErrNotExist = errors.New("not exist")
)
