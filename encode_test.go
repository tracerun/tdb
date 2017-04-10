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
	fRange := newFileRange(0, 0)
	assert.Equal(t, fileEncode(0), fRange.start, "start should be 0")
	assert.Equal(t, fileEncode(math.MaxUint32), fRange.end, "end should be the max uint32")
}
