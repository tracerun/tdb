package tdb

import (
	"os"
)

// TDB is the entrance point.
type TDB struct {
	path string // the folder path used for tdb

	slot *info // slot index info
	meta *info // meta info
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

	// load slot index
	if err := db.loadSlotIndex(); err != nil {
		return nil, err
	}

	return db, nil
}
