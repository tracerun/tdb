package tdb

import (
	"math/rand"
	"path"
)

const (
	slotIndex   = "__slotindex__"
	letterBytes = "123456789abcdefghijklmnopqrstuvwxyz0"
	letterLen   = 36
)

func randBytes(n int) []byte {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(letterLen)]
	}
	return b
}

func (db *TDB) loadSlotIndex() error {
	var err error
	db.slot, err = createInfo(path.Join(db.path, slotIndex))
	return err
}

// AddSlot to add a slot to database
func (db *TDB) AddSlot(target string, start, howlong uint32) {

}
