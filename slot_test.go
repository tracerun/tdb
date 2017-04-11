package tdb

import (
	"os"
	"path/filepath"
	"strconv"
	"testing"
	"time"

	"fmt"

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

	var thisDay string
	thisOffset := time.Date(2017, 4, day, 0, 0, 0, 0, time.Local).Unix()
	if hours >= 12 {
		thisDay = strconv.Itoa(day*10 + 5)
		thisOffset = thisOffset + 43200
	} else {
		thisDay = strconv.Itoa(day * 10)
	}

	if len(thisDay) == 2 {
		thisDay = fmt.Sprintf("0%s", thisDay)
	}

	assert.Equal(t, "201704", folder, "folder is wrong")
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

	// create slot1 and add it
	start1 := uint32(1491134201) // 201704020 (UTC)
	slot1 := uint32(20)          // 20 seconds
	err = db.AddSlot(target, start1, slot1)
	assert.NoError(t, err, "should have no error")

	// create slot2 and add it
	start2 := uint32(1551134201) // 201902255 (UTC)
	slot2 := uint32(40)          // 40 seconds
	err = db.AddSlot(target, start2, slot2)
	assert.NoError(t, err, "should have no error")

	targets := db.GetTargets()
	assert.Len(t, targets, 1, "should have one target.")

	files, err := db.getTargetFiles(target, 0, 0)
	assert.NoError(t, err, "should have no error while get target files")
	assert.Len(t, files, 2, "should have 2 slots")
}
