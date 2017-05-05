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
	// ErrRange the range is wrong
	ErrRange = errors.New("range is wrong")
	// ErrActionValue error for action value field
	ErrActionValue = errors.New("action value bytes wrong")
	// ErrProjectUnavailable error for project handle
	ErrProjectUnavailable = errors.New("project not available")
	// ErrTargetNotBelongToProject error for wrong target to project
	ErrTargetNotBelongToProject = errors.New("target not belong to project")
)
