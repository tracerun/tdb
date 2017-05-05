package tdb

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

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
	targetHome, err := db.getAliasedHome(target)
	assert.NoError(t, err, "fail to get target home folder")

	exist, err := checkFolderExist(targetHome)
	assert.NoError(t, err, "target home folder wrong")
	assert.True(t, exist, "target home folder should exist")

	b := db.slot.getInfoValue(target)
	assert.Len(t, string(b), slotAliasLen, "alias name length wrong")

	thisTargetHome := filepath.Join(db.path, slotsFolder, string(b))
	assert.Equal(t, thisTargetHome, targetHome, "target home not correct")

	// create slot1 and add it
	start1 := uint32(1491134201) // 201704020 (UTC)
	slot1 := uint32(20)          // 20 seconds
	err = db.addSlot(target, start1, slot1)
	assert.NoError(t, err, "should have no error")

	// create slot2 and add it
	start2 := uint32(1551134201) // 201902255 (UTC)
	slot2 := uint32(40)          // 40 seconds
	err = db.addSlot(target, start2, slot2)
	assert.NoError(t, err, "should have no error")

	targets := db.GetTargets()
	assert.Len(t, targets, 1, "should have one target.")

	files, err := db.getTargetFiles(target, 0, 0)
	assert.NoError(t, err, "should have no error while get target files")
	assert.Len(t, files, 2, "should have 2 slots")

	// create slot3 and add it
	start3 := uint32(1491134202) // 201704020 (UTC)
	slot3 := uint32(10)          // 10 seconds
	err = db.addSlot(target, start3, slot3)
	assert.NoError(t, err, "should have no error")

	// get slots with bad range
	_, _, err = db.GetSlots(target, 20, 10)
	assert.Equal(t, ErrRange, err, "range should be wrong")

	// get one
	starts, slots, err := db.GetSlots(target, 1491134203, 0)
	assert.NoError(t, err, "should have no error to get all slots")
	assert.Len(t, starts, 2, "should have two files")
	assert.Len(t, slots, 2, "should have two files")

	assert.Len(t, starts[0], 0, "should have no slot")
	assert.Len(t, slots[0], 0, "should have no slot")
	assert.Len(t, starts[1], 1, "should have one slot")
	assert.Len(t, slots[1], 1, "should have one slot")

	// get all slots
	starts, slots, err = db.GetSlots(target, 0, 0)
	assert.NoError(t, err, "should have no error to get all slots")
	assert.Len(t, starts, 2, "should have two files")
	assert.Len(t, slots, 2, "should have two files")

	assert.Equal(t, start1, starts[0][0], "start1 wrong")
	assert.Equal(t, slot1, slots[0][0], "slot1 wrong")
	assert.Equal(t, start3, starts[0][1], "start3 wrong")
	assert.Equal(t, slot3, slots[0][1], "slot3 wrong")
	assert.Equal(t, start2, starts[1][0], "start2 wrong")
	assert.Equal(t, slot2, slots[1][0], "slot2 wrong")
}
