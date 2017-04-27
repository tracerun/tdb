package tdb

import (
	"path/filepath"
	"time"

	"go.uber.org/zap"
)

const (
	actionIndex = "__action__"

	actionExp = 15 // action will expire after 15 seconds
)

// AddAction used to add actions to database
func (db *TDB) AddAction(target string, ts uint32) error {
	db.action.contentLock.Lock()

	err := handleAction(db, target, ts)
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

// CheckExpirations to check expired actions
func (db *TDB) CheckExpirations() error {
	targets, starts, lasts, err := db.GetActions()
	if err != nil {
		return err
	}

	changed := false
	now := uint32(time.Now().Unix())

	db.action.contentLock.Lock()
	for i := 0; i < len(lasts); i++ {
		if now-lasts[i] > actionExp {
			changed = true
			p("expired action", zap.String("target", targets[i]))
			if err := db.AddSlot(targets[i], starts[i], lasts[i]-starts[i]); err != nil {
				db.action.contentLock.Unlock()
				return err
			}
			delete(db.action.content, targets[i])
		}
	}
	db.action.contentLock.Unlock()

	if changed {
		return db.action.writeToDisk()
	}
	return nil
}

func handleAction(db *TDB, target string, ts uint32) error {
	actions := db.action.content
	v, ok := actions[target]
	if !ok {
		p("new action")
		actions[target] = encodeAction(ts, ts+1)
		return nil
	}
	start, last, err := decodeAction(v)
	if err != nil {
		return err
	}

	if ts > last {
		if ts-last > actionExp {
			p("new slot", zap.Uint32("ts", ts), zap.Uint32("last", last))
			if err := db.AddSlot(target, start, last-start); err != nil {
				return err
			}
			actions[target] = encodeAction(ts, ts+1)
		} else {
			actions[target] = encodeAction(start, ts)
		}
	}
	return nil
}

// load slot index from file
func (db *TDB) loadActionIndex() error {
	var err error
	db.action, err = createInfo(filepath.Join(db.path, actionIndex))
	return err
}
