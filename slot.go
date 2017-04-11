package tdb

import (
	"encoding/binary"
	"math/rand"
	"os"
	"path/filepath"
	"sort"
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

	offsetExt = ".idx"
	slotExt   = ".slt"
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

	file := encodeFileFromUnix(start)
	offset := start - file.origin()
	return writeSlotToFile(targetHome, file, uint16(offset), howlong)
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

// the file encode used currently
func currentFileEncode() fileEncode {
	now := time.Now().Unix()
	return encodeFileFromUnix(uint32(now))
}

func writeSlotToFile(tagHome string, file fileEncode, offset uint16, howlong uint32) error {
	path, fileName := file.path()

	folder := filepath.Join(tagHome, path)
	// create folder if not exist
	if err := createFolder(folder); err != nil {
		return err
	}

	// append to offset file
	offsetFile := strings.Join([]string{fileName, offsetExt}, "")
	offsetB := make([]byte, 2)
	binary.LittleEndian.PutUint16(offsetB, offset)
	if err := appendToFile(filepath.Join(folder, offsetFile), offsetB); err != nil {
		return err
	}

	// append to slot file
	slotFile := strings.Join([]string{fileName, slotExt}, "")
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
func (db *TDB) getTargetFiles(target string, start, end uint32) ([]fileEncode, error) {
	tagHome, err := db.getTargetHome(target)
	if err != nil {
		return nil, err
	}

	// create the range
	fRange, err := newFileRange(start, end)
	if err != nil {
		return nil, err
	}

	// get folders within the range
	folders, err := getInRangeFolders(fRange, tagHome)
	if err != nil {
		return nil, err
	}

	var fileEncodes []fileEncode
	for _, folder := range folders {
		files, err := getInRangeFiles(fRange, tagHome, folder)
		if err != nil {
			return nil, err
		}
		fileEncodes = append(fileEncodes, files...)
	}

	return fileEncodes, nil
}

func getInRangeFolders(fRange *fileRange, targetHome string) ([]string, error) {
	f, err := os.Open(targetHome)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var fs []string

	names, err := f.Readdirnames(0)
	if err != nil {
		return nil, err
	}

	for _, n := range names {
		// check whether this folder is in range
		in, err := fRange.folderInRange(n)
		if err != nil {
			return nil, err
		} else if in {
			fs = append(fs, n)
		}
	}

	// sort to ensure the strings in increasing order
	if !sort.StringsAreSorted(fs) {
		sort.Strings(fs)
	}

	return fs, err
}

// get sorted []fileEncode
func getInRangeFiles(fRange *fileRange, path, folder string) ([]fileEncode, error) {
	f, err := os.Open(filepath.Join(path, folder))
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var filesIn []fileEncode
	names, err := f.Readdirnames(0)
	if err != nil {
		return nil, err
	}

	for _, n := range names {
		// check whether this file is in range
		if filepath.Ext(n) != offsetExt {
			continue
		}

		encoded, err := encodeFromPath(folder, n)
		if err != nil {
			return nil, err
		}

		in, err := fRange.fileInRange(encoded)
		if err != nil {
			return nil, err
		}

		if in {
			filesIn = append(filesIn, encoded)
		}
	}

	s := fileEncodeSlice(filesIn)
	if !sort.IsSorted(s) {
		sort.Sort(s)
	}

	return filesIn, nil
}
