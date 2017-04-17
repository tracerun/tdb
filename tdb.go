package tdb

// TDB is the entrance point.
type TDB struct {
	path string // the folder path used for tdb

	slot   *info // slot index info
	meta   *info // meta info
	action *info // record actions
}

// Open creates an instance of TDB.
// "p" should be a path for folder, if not exist, create one.
func Open(p string) (*TDB, error) {
	db := &TDB{path: p}
	if err := createFolder(p); err != nil {
		return nil, err
	}

	// load meta information
	if err := db.loadMeta(); err != nil {
		return nil, err
	}

	// load slot index
	if err := db.loadSlotIndex(); err != nil {
		return nil, err
	}

	// load action index
	if err := db.loadActionIndex(); err != nil {
		return nil, err
	}

	return db, nil
}
