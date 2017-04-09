package tdb

import (
	"testing"

	"path/filepath"

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
	assert.Equal(t, filepath.Join("2017", "12"), folder, "folder is wrong")
	assert.Equal(t, "1p", filename, "folder is wrong")
}
