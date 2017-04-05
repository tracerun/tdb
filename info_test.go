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

	// creation on folder will fail
	_, err := createInfo(folderPath)
	assert.Equal(t, ErrInfoFilePath, err, "error should pop while creating info on a folder")

	demoInfo := "demo_info"
	defer os.Remove(demoInfo)
	one, err := createInfo(demoInfo)
	assert.NoError(t, err, "should have no error while creating info from non-existed file")
	assert.Len(t, one.content, 0, "should have empty content")

	key := "key"
	value := []byte("value")
	one.setContent([]string{key}, [][]byte{value})
	assert.Len(t, one.content, 1, "should have one content")

	two, err2 := createInfo(demoInfo)
	assert.NoError(t, err2, "should have no error while creating info from an existed file")
	assert.Len(t, two.content, 1, "should have empty content")
	assert.EqualValues(t, one.content[key], two.content[key], "the values should be the same")
}
