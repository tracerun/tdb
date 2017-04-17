package tdb

import (
	"path/filepath"
)

const (
	actionIndex = "__action__"

	actionExp = 15 // action will expire after 15 seconds
)

// AddAction used to add actions to database
func (db *TDB) AddAction(target string, active bool, ts uint32) error {
	db.action.contentLock.Lock()

	err := handleAction(db, target, active, ts)
	if err != nil {
		db.action.contentLock.Unlock()
		return err
	}
	db.action.contentLock.Unlock()

	return db.action.writeToDisk()
}

// GetActions to get all actions.
// return targets, starts, lasts, error
func (db *TDB) GetActions() ([]string, []uint32, []uint32, error) {
	keys, values := db.action.getAllInfo()

	var targets []string
	var starts, lasts []uint32

	for i := 0; i < len(keys); i++ {
		start, last, err := decodeAction(values[i])
		if err != nil {
			return nil, nil, nil, err
		}
		targets = append(targets, keys[i])
		starts = append(starts, start)
		lasts = append(lasts, last)
	}
	return targets, starts, lasts, nil
}

func handleAction(db *TDB, target string, active bool, ts uint32) error {
	actions := db.action.content
	v := actions[target]
	if v == nil {
		if active {
			actions[target] = encodeAction(ts, ts+1)
			return nil
		}
		return db.AddSlot(target, ts, uint32(1))
	}
	start, last, err := decodeAction(v)
	if err != nil {
		return err
	}

	if active {
		if ts-last > actionExp {
			if err := db.AddSlot(target, start, last-start); err != nil {
				return err
			}
			actions[target] = encodeAction(ts, ts+1)
		} else if ts > last {
			actions[target] = encodeAction(start, ts)
		}
	} else {
		if ts-last > actionExp {
			if err := db.AddSlot(target, start, last-start); err != nil {
				return err
			}
			if err := db.AddSlot(target, ts, uint32(1)); err != nil {
				return err
			}
		} else {
			howlong := ts - start
			if ts < last {
				howlong = last - start
			}
			if err := db.AddSlot(target, start, howlong); err != nil {
				return err
			}
		}
		delete(actions, target)
	}
	return nil
}

// load slot index from file
func (db *TDB) loadActionIndex() error {
	var err error
	db.action, err = createInfo(filepath.Join(db.path, actionIndex))
	return err
}
