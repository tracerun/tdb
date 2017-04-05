package tdb

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInfoMethods(t *testing.T) {
	// create a bad folder
	folderPath := "bad_folder"

	if _, err := os.Stat(folderPath); err != nil {
		if os.IsNotExist(err) {
			if err := os.MkdirAll(folderPath, os.ModePerm); err != nil {
				t.Error(err)
			}
		} else {
			t.Fatal(err)
		}
	}
	defer os.RemoveAll(folderPath)

	_, err := createInfo(folderPath)
	assert.Equal(t, ErrInfoFilePath, err, "error should pop while creating info on a folder")

	demoInfo := "demo_info"
	one, err := createInfo(demoInfo)
	assert.NoError(t, err, "should have no error while creating info from non-existed file")
	assert.Len(t, one.content, 0, "should have empty content")

	one.setContent([]string{"key"}, [][]byte{[]byte("value")})
	assert.Len(t, one.content, 1, "should have empty content")
}
