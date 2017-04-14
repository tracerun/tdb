package tdb

import (
	"math"
	"sort"
	"testing"

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
	assert.Equal(t, "015", filename, "folder is wrong")
}

func TestFileEncodeFromPath(t *testing.T) {
	encoded, err := encodeFromPath("201705", "025.idx")
	assert.NoError(t, err, "should have no error when encoding from path")
	assert.Equal(t, fileEncode(201705025), encoded, "file encode wrong")

	_, err = encodeFromPath("201705a", "025.idx")
	assert.NotNil(t, err, "should have error when folder name length is not 6")

	_, err = encodeFromPath("201705", "0025.idx")
	assert.NotNil(t, err, "should have error when file name length is not 3")

	_, err = encodeFromPath("20170a", "025.idx")
	assert.NotNil(t, err, "bad folder name when encoding")

	_, err = encodeFromPath("201705", "a25.idx")
	assert.NotNil(t, err, "bad file name when encoding")
}

func TestSortingFileEncode(t *testing.T) {
	var fileEncodes []fileEncode
	fileEncodes = append(fileEncodes, fileEncode(3))
	fileEncodes = append(fileEncodes, fileEncode(2))
	fileEncodes = append(fileEncodes, fileEncode(1))

	s := fileEncodeSlice(fileEncodes)
	sort.Sort(s)

	assert.Equal(t, fileEncode(1), fileEncodes[0], "first value wrong")
	assert.Equal(t, fileEncode(2), fileEncodes[1], "second value wrong")
	assert.Equal(t, fileEncode(3), fileEncodes[2], "third value wrong")
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

	// from 201704020 to 201902255 (UTC)
	fRange, err = newFileRange(1491134201, 1551134201)

	// test folder range
	inRange, err = fRange.folderInRange("201712")
	assert.True(t, inRange, "should in range")
	assert.NoError(t, err, "should have no error")

	inRange, err = fRange.folderInRange("201903")
	assert.False(t, inRange, "should not in range")
	assert.NoError(t, err, "should have no error")

	inRange, err = fRange.folderInRange("201705120")
	assert.False(t, inRange, "should not in range")
	assert.NoError(t, err, "should have no error")

	// test file range
	inRange, err = fRange.fileInRange(fileEncode(201801250))
	assert.NoError(t, err, "should have no error")
	assert.True(t, inRange, "should in range")

	inRange, err = fRange.fileInRange(fileEncode(202001250))
	assert.NoError(t, err, "should have no error")
	assert.False(t, inRange, "should in range")
}

func TestEncodeAliasAndFile(t *testing.T) {
	alias := "lkfj/abcdef"
	file := fileEncode(201501205)
	assert.Equal(t, "abcdef201501205", encodeAliasAndFile(alias, file), "encode wrong")
}
