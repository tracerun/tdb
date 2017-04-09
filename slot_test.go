package tdb

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRandBytes(t *testing.T) {
	assert.Len(t, randBytes(10), 10, "bytes length should be 10")
}

func BenchmarkRandBytes(b *testing.B) {
	for i := 0; i < b.N; i++ {
		randBytes(10)
	}
}

func TestGetDetailFile(t *testing.T) {
	folder, file, offset := getDetailFile(uint32(1491550758))

	thisTime := time.Unix(int64(1491550758), 0)
	_, _, day := thisTime.Date()
	hours := thisTime.Hour()

	thisDay := strconv.Itoa(day)

	thisOffset := time.Date(2017, 4, day, 0, 0, 0, 0, time.Local).Unix()
	if hours >= 12 {
		thisDay = fmt.Sprintf("%sp", thisDay)
		thisOffset = thisOffset + 43200
	}

	thisPath := filepath.Join("2017", "4")
	assert.Equal(t, thisPath, folder, "folder is wrong")
	assert.Equal(t, thisDay, file, "file is wrong")
	assert.Equal(t, uint16(1491550758-thisOffset), offset, "offset is wrong")
}

func TestWriteSlotToFile(t *testing.T) {
	folder := "test_write_slot"
	file := "test"
	defer os.RemoveAll(folder)

	err := writeSlotToFile(folder, file, uint16(123), uint32(123))
	assert.NoError(t, err, "append to file wrong")
}

func TestSlots(t *testing.T) {
	demoSlot := "slot_test"
	db, err := Open(demoSlot)
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(db.path)

	target := "abc.edf"
	targetHome, err := db.getTargetHome(target)
	assert.NoError(t, err, "fail to get target home folder")

	exist, err := checkFolderExist(targetHome)
	assert.NoError(t, err, "target home folder wrong")
	assert.True(t, exist, "target home folder should exist")

	b := db.slot.getValue(target)
	assert.Len(t, string(b), slotAliasLen, "alias name length wrong")

	thisTargetHome := filepath.Join(db.path, slotsFolder, string(b))
	assert.Equal(t, thisTargetHome, targetHome, "target home not correct")
}
