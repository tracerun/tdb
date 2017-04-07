package tdb

import (
	"path/filepath"
	"testing"

	"os"

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

	thisPath := filepath.Join("2017", "4")
	assert.Equal(t, thisPath, folder, "folder is wrong")
	assert.Equal(t, "7p", file, "file is wrong")
	assert.Equal(t, uint16(13158), offset, "offset is wrong")

	folder, file, offset = getDetailFile(uint32(1491532758))
	assert.Equal(t, thisPath, folder, "folder is wrong")
	assert.Equal(t, "7", file, "file is wrong")
	assert.Equal(t, uint16(38358), offset, "offset is wrong")
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

	stat, err := os.Stat(targetHome)
	assert.NoError(t, err, "target home folder wrong")

	if !stat.IsDir() {
		t.Error("target home not a dir")
	}

	b := db.slot.getValue(target)
	assert.Len(t, string(b), slotAliasLen, "alias name length wrong")

	thisTargetHome := filepath.Join(db.path, slotsFolder, string(b))
	assert.Equal(t, thisTargetHome, targetHome, "target home not correct")
}
