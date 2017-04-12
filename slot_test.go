package tdb

import (
	"os"
	"path/filepath"
	"testing"

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

func TestWriteSlotToFile(t *testing.T) {
	folder := "test_write_slot"
	file := encodeFile(2017, 4, 11, 23)
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
