package tdb

import (
	"testing"

	"math"

	"github.com/stretchr/testify/assert"
)

func TestFileEncode(t *testing.T) {
	encoded := encodeFile(2017, 4, 23, 5)
	assert.Equal(t, fileEncode(201704230), encoded, "file encode wrong")
	assert.True(t, encoded.isAM(), "should be in am")

	encoded = encodeFile(2017, 12, 1, 23)
	assert.Equal(t, fileEncode(201712015), encoded, "file encode wrong")

	assert.Equal(t, 2017, encoded.year(), "encode year wrong")
	assert.Equal(t, 12, encoded.month(), "encode month wrong")
	assert.Equal(t, 1, encoded.day(), "encode month wrong")
	assert.False(t, encoded.isAM(), "should be in pm")

	folder, filename := encoded.path()
	assert.Equal(t, "201712", folder, "folder is wrong")
	assert.Equal(t, "15", filename, "folder is wrong")
}

func TestFileRange(t *testing.T) {
	fRange, err := newFileRange(0, 0)
	assert.NoError(t, err, "error creating file range")
	assert.Equal(t, fileEncode(0), fRange.start, "start should be 0")
	assert.Equal(t, fileEncode(math.MaxUint32), fRange.end, "end should be the max uint32")

	fRange, err = newFileRange(2, 1)
	assert.Equal(t, ErrRange, err, "should have error when start > end")
	assert.Nil(t, fRange, "filerange should be nil")

	fRange, err = newFileRange(0, 0)
	inRange, err := fRange.folderInRange("badfolder")
	assert.False(t, inRange, "should not in range.")
	assert.NotNil(t, err, "can't convert to int")

	// from 201704 to 201902
	fRange, err = newFileRange(1491134201, 1551134201)

	inRange, err = fRange.folderInRange("201712")
	assert.True(t, inRange, "should in range")
	assert.NoError(t, err, "should have no error")

	inRange, err = fRange.folderInRange("201903")
	assert.False(t, inRange, "should not in range")
	assert.NoError(t, err, "should have no error")

	inRange, err = fRange.folderInRange("201705120")
	assert.False(t, inRange, "should not in range")
	assert.NoError(t, err, "should have no error")
}
