package tdb

import (
	"encoding/binary"
	"io"
	"math"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"go.uber.org/zap"

	"github.com/tracerun/locker"
	"golang.org/x/sync/errgroup"
)

const (
	slotIndex = "__slotindex__"

	slotAliasLen = 6
	slotsFolder  = "slots"
	slotBytes    = 4
	indexBytes   = 2

	offsetExt = ".idx"
	slotExt   = ".slt"
)

var (
	srcLocker *locker.Locker
)

func init() {
	srcLocker = locker.New()
}

// AddSlot to add a slot to database
func (db *TDB) AddSlot(target string, start, howlong uint32) error {
	aliasedHome, err := db.getAliasedHome(target)
	if err != nil {
		return err
	}

	file := encodeFileFromUnix(start)
	offset := start - file.origin()
	return writeSlotToFile(aliasedHome, file, uint16(offset), howlong)
}

// GetTargets to get all the targets
func (db *TDB) GetTargets() []string {
	return db.slot.getInfoKeys()
}

// GetSlots to get slots of a target in certain range
// return unix time and slots
func (db *TDB) GetSlots(target string, start, end uint32) ([][]uint32, [][]uint32, error) {
	aliasName := string(db.slot.getInfoValue(target))
	if aliasName == "" {
		return nil, nil, nil
	}

	files, err := db.getTargetFiles(target, start, end)
	if err != nil {
		return nil, nil, err
	}

	aliasedHome := filepath.Join(db.path, slotsFolder, aliasName)

	var g errgroup.Group
	starts := make([][]uint32, len(files))
	slots := make([][]uint32, len(files))

	realStart := start
	readEnd := end
	if readEnd == 0 {
		readEnd = math.MaxUint32
	}

	for i, f := range files {
		i, f := i, f
		g.Go(func() error {
			var startsResult, slotsResult []uint32
			thisStarts, thisSlots, err := readFile(aliasedHome, f)
			p("one file", zap.String("alias", aliasName), zap.Any("starts", thisStarts))
			// get the in range results
			for i := 0; i < len(thisStarts); i++ {
				if thisStarts[i] >= realStart && thisStarts[i] <= readEnd {
					startsResult = append(startsResult, thisStarts[i])
					slotsResult = append(slotsResult, thisSlots[i])
				}
			}
			starts[i], slots[i] = startsResult, slotsResult
			return err
		})
	}
	err = g.Wait()
	return starts, slots, err
}

func readFile(aliasedHome string, file fileEncode) ([]uint32, []uint32, error) {
	path, fileName := file.path()
	baseName := filepath.Join(aliasedHome, path, fileName)

	idxName := strings.Join([]string{baseName, offsetExt}, "")
	slotName := strings.Join([]string{baseName, slotExt}, "")

	unlock := srcLocker.ReadLock(encodeAliasAndFile(aliasedHome, file))
	defer unlock()

	// open offset index file
	idxFile, err := os.Open(idxName)
	if err != nil {
		return nil, nil, err
	}
	defer idxFile.Close()

	// open slot file
	slotFile, err := os.Open(slotName)
	if err != nil {
		return nil, nil, err
	}
	defer slotFile.Close()

	offsetB := make([]byte, 2)
	slotB := make([]byte, 4)

	origin := file.origin()
	var starts, slots []uint32
	for {
		// read offset with 2 bytes
		n, err := idxFile.Read(offsetB)
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, nil, err
		}
		if n != 2 {
			return nil, nil, err
		}
		offset := binary.LittleEndian.Uint16(offsetB)

		// read slot with 4 bytes
		n, err = slotFile.Read(slotB)
		if err != nil {
			return nil, nil, err
		}
		if n != 4 {
			return nil, nil, err
		}
		slot := binary.LittleEndian.Uint32(slotB)

		starts = append(starts, origin+uint32(offset))
		slots = append(slots, slot)
	}
	return starts, slots, nil
}

// load slot index from file
func (db *TDB) loadSlotIndex() error {
	var err error
	db.slot, err = createInfo(filepath.Join(db.path, slotIndex))
	return err
}

// get aliased home folder for the target
// If target is not exist, create it in the index and also create a folder for it.
func (db *TDB) getAliasedHome(target string) (string, error) {
	b := db.slot.getInfoValue(target)
	if b != nil {
		return filepath.Join(db.path, slotsFolder, string(b)), nil
	}

	// not exist, create folder
	unlock := srcLocker.WriteLock(target)
	defer unlock()

	b = db.slot.getInfoValue(target)
	if b != nil {
		return filepath.Join(db.path, slotsFolder, string(b)), nil
	}

	var aliasedHome string
	var aliasName string
	for {
		aliasName = string(randBytes(slotAliasLen))
		aliasedHome = filepath.Join(db.path, slotsFolder, aliasName)

		// check whether this folder is exist
		exist, err := checkFolderExist(aliasedHome)
		if err != nil {
			return "", err
		}

		// no exist, OK, create it.
		if !exist {
			err := createFolder(aliasedHome)
			if err != nil {
				return "", err
			}
			// successfully created the folder, BREAK
			break
		}
	}
	return aliasedHome, db.slot.updateInfo([]string{target}, [][]byte{[]byte(aliasName)})
}

func writeSlotToFile(aliasedHome string, file fileEncode, offset uint16, howlong uint32) error {
	subFolder, fileName := file.path()

	fullFolder := filepath.Join(aliasedHome, subFolder)
	// create folder if not exist
	if err := createFolder(fullFolder); err != nil {
		return err
	}

	unlock := srcLocker.WriteLock(encodeAliasAndFile(aliasedHome, file))
	defer unlock()

	// append to offset file
	offsetFileName := strings.Join([]string{fileName, offsetExt}, "")
	offsetB := make([]byte, 2)
	binary.LittleEndian.PutUint16(offsetB, offset)
	if err := appendToFile(filepath.Join(fullFolder, offsetFileName), offsetB); err != nil {
		return err
	}

	// append to slot file
	slotFileName := strings.Join([]string{fileName, slotExt}, "")
	slotB := make([]byte, 4)
	binary.LittleEndian.PutUint32(slotB, howlong)
	return appendToFile(filepath.Join(fullFolder, slotFileName), slotB)
}

func appendToFile(fullPath string, b []byte) error {
	f, err := os.OpenFile(fullPath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)
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
	aliasedHome, err := db.getAliasedHome(target)
	if err != nil {
		return nil, err
	}

	// create the range
	fRange, err := newFileRange(start, end)
	if err != nil {
		return nil, err
	}

	// get folders within the range
	folders, err := getInRangeFolders(fRange, aliasedHome)
	if err != nil {
		return nil, err
	}

	var fileEncodes []fileEncode
	for _, folder := range folders {
		files, err := getInRangeFiles(fRange, aliasedHome, folder)
		if err != nil {
			return nil, err
		}
		fileEncodes = append(fileEncodes, files...)
	}

	return fileEncodes, nil
}

func getInRangeFolders(fRange *fileRange, aliasedHome string) ([]string, error) {
	f, err := os.Open(aliasedHome)
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
func getInRangeFiles(fRange *fileRange, aliasedHome, subFolder string) ([]fileEncode, error) {
	f, err := os.Open(filepath.Join(aliasedHome, subFolder))
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

		encoded, err := encodeFromPath(subFolder, n)
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
