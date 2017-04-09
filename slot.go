package tdb

import (
	"encoding/binary"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"sync"
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
	fileLocker *slotFileLocker
)

type slotFileLocker struct {
	files map[string]*sync.RWMutex
	lock  *sync.RWMutex
}

func init() {
	slotLocker = new(sync.RWMutex)

	fileLocker = &slotFileLocker{
		files: make(map[string]*sync.RWMutex),
		lock:  new(sync.RWMutex),
	}
}

// AddSlot to add a slot to database
func (db *TDB) AddSlot(target string, start, howlong uint32) error {
	targetHome, err := db.getTargetHome(target)
	if err != nil {
		return err
	}

	folder, fileName, offset := getDetailFile(start)
	fileFolder := filepath.Join(targetHome, folder)

	return writeSlotToFile(fileFolder, fileName, offset, howlong)
}

// GetTargets to get all the targets
func (db *TDB) GetTargets() []string {
	return db.slot.getKeys()
}

// GetSlots to get all the slots for a target
// return unix time and slots
func (db *TDB) GetSlots(target string) (starts []uint32, slots []uint32) {
	aliasName := string(db.slot.getValue(target))
	if aliasName == "" {
		return
	}
	return
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

		// check whether this folder is exist
		exist, err := checkFolderExist(targetHome)
		if err != nil {
			return "", err
		}

		// no exist, OK, create it.
		if !exist {
			err := createFolder(targetHome)
			if err != nil {
				return "", err
			}
			// successfully created the folder, BREAK
			break
		}
	}
	return targetHome, db.slot.updateInfo([]string{target}, [][]byte{[]byte(aliasName)})
}

// start is a unixtime
// return which folder should be write to, filename and offset
func getDetailFile(start uint32) (string, string, uint16) {
	encodedStart := encodeFileFromUnix(start)

	folder, fileName := encodedStart.path()
	offset := start - encodedStart.origin()

	return folder, fileName, uint16(offset)
}

func writeSlotToFile(folder string, file string, offset uint16, howlong uint32) error {
	// create folder if not exist
	if err := createFolder(folder); err != nil {
		return err
	}

	// append to offset file
	offsetFile := strings.Join([]string{file, offsetExt}, ".")
	offsetB := make([]byte, 2)
	binary.LittleEndian.PutUint16(offsetB, offset)
	if err := appendToFile(filepath.Join(folder, offsetFile), offsetB); err != nil {
		return err
	}

	// append to slot file
	slotFile := strings.Join([]string{file, slotExt}, ".")
	slotB := make([]byte, 4)
	binary.LittleEndian.PutUint32(slotB, howlong)
	return appendToFile(filepath.Join(folder, slotFile), slotB)
}

func appendToFile(fileName string, b []byte) error {
	f, err := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write(b)
	return err
}

// Only get target files that contain slots between start and end.
// start, end should be unixtime
func (db *TDB) getTargetFiles(target string, start, end uint32) ([]uint32, []string, error) {
	if start != 0 && end != 0 && start > end {
		return nil, nil, ErrRange
	}

	// tagHome, err := db.getTargetHome(target)
	// if err != nil {
	// 	return nil, nil, err
	// }

	var starts []uint32
	var files []string

	// if start == 0 {
	// 	startTime = 0
	// } else {
	// 	startT := time.Unix(int64(start), 0)

	// }

	// if end == 0 {
	// 	endTime = math.MaxFloat32
	// } else {

	// }
	return starts, files, nil
}
