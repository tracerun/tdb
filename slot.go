package tdb

import (
	"encoding/binary"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	slotIndex   = "__slotindex__"
	letterBytes = "0123456789abcdefghijklmnopqrstuvwxyz"
	letterLen   = 36

	slotAliasLen = 6
	slotsFolder  = "slots"
	slotBytes    = 4
	indexBytes   = 2

	offsetExt = "idx"
	slotExt   = "slt"
)

var (
	slotLocker *sync.RWMutex
)

func init() {
	slotLocker = new(sync.RWMutex)
}

func randBytes(n int) []byte {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(letterLen)]
	}
	return b
}

// load slot index from file
func (db *TDB) loadSlotIndex() error {
	var err error
	db.slot, err = createInfo(filepath.Join(db.path, slotIndex))
	return err
}

// get home folder for the target
// If target is not exist, create it in the index and also create a folder for it.
func (db *TDB) getTargetHome(target string) (string, error) {
	slotLocker.RLock()
	b := db.slot.getValue(target)
	slotLocker.RUnlock()
	if b != nil {
		return filepath.Join(db.path, slotsFolder, string(b)), nil
	}

	// not exist, create folder
	slotLocker.Lock()
	defer slotLocker.Unlock()

	var targetHome string
	var aliasName string
	for {
		aliasName = string(randBytes(slotAliasLen))
		targetHome = filepath.Join(db.path, slotsFolder, aliasName)

		if stat, err := os.Stat(targetHome); err != nil {
			if os.IsNotExist(err) {
				// after successfully create the folder, break
				if err := os.MkdirAll(targetHome, os.ModePerm); err != nil {
					return "", err
				}
				break
			} else {
				return "", err
			}
		} else if !stat.IsDir() {
			return "", ErrDBPathNotFolder
		}
	}
	return targetHome, db.slot.updateInfo([]string{target}, [][]byte{[]byte(aliasName)})
}

// start is a unixtime
// return which folder should be write to, filename and offset
func getDetailFile(start uint32) (string, string, uint16) {
	t := time.Unix(int64(start), 0)
	year, month, day := t.Date()
	hours := t.Hour()

	fileName := strconv.Itoa(day)
	fileOrigin := time.Date(year, month, day, 0, 0, 0, 0, time.Local).Unix()
	if hours >= 12 {
		fileOrigin = fileOrigin + 43200
		fileName = fmt.Sprintf("%sp", fileName)
	}

	offset := start - uint32(fileOrigin)
	return filepath.Join(strconv.Itoa(year), strconv.Itoa(int(month))), fileName, uint16(offset)
}

func writeSlotToFile(file string, offset uint16, howlong uint32) error {
	// append to offset file
	offsetFile := strings.Join([]string{file, offsetExt}, ".")
	offsetF, err := os.OpenFile(offsetFile, os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer offsetF.Close()

	var b []byte
	binary.LittleEndian.PutUint16(b, offset)
	if _, err := offsetF.Write(b); err != nil {
		return err
	}

	// append to slot file
	slotFile := strings.Join([]string{file, slotExt}, ".")
	slotF, err := os.OpenFile(slotFile, os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer slotF.Close()

	binary.LittleEndian.PutUint32(b, howlong)
	if _, err := slotF.Write(b); err != nil {
		return err
	}

	return nil
}

// AddSlot to add a slot to database
func (db *TDB) AddSlot(target string, start, howlong uint32) error {
	targetHome, err := db.getTargetHome(target)
	if err != nil {
		return err
	}

	folder, fileName, offset := getDetailFile(start)
	file := filepath.Join(targetHome, folder, fileName)

	return writeSlotToFile(file, offset, howlong)
}
