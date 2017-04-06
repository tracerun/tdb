package tdb

import (
	"os"
)

// TDB is the entrance point.
type TDB struct {
	// the folder path used for tdb
	path string

	slot *info
	meta *info
}

// Open creates an instance of TDB.
// "p" should be a path for folder, if not exist, create one.
func Open(p string) (*TDB, error) {
	db := &TDB{path: p}
	if stat, err := os.Stat(p); err != nil {
		if os.IsNotExist(err) {
			if err := os.MkdirAll(p, os.ModePerm); err != nil {
				return nil, err
			}
		} else {
			return db, err
		}
	} else if !stat.IsDir() {
		return db, ErrDBPathNotFolder
	}

	// load meta information
	if err := db.loadMeta(); err != nil {
		return nil, err
	}

	return db, nil
}
