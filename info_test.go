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
	err = one.update([]string{key}, [][]byte{value})
	assert.NoError(t, err, "should have no error while updating the content")
	assert.Len(t, one.content, 1, "should have one content")

	// k, v length not equal
	err = one.update([]string{key}, [][]byte{value, value})
	assert.Error(t, err, "should have error when k, v length not equal")

	// empty k, v
	err = one.update([]string{}, [][]byte{})
	assert.NoError(t, err, "should have no error while updating nothing")

	// load an existing content
	two, err2 := createInfo(demoInfo)
	assert.NoError(t, err2, "should have no error while creating info from an existed file")
	assert.Len(t, two.content, 1, "should have empty content")
	assert.EqualValues(t, one.content[key], two.content[key], "the values should be the same")
}

func BenchmarkLoad(b *testing.B) {
	benchInfo := "bench_info_load"
	defer os.Remove(benchInfo)
	one, err := createInfo(benchInfo)
	if err != nil {
		b.Fatal(err)
	}
	if len(one.content) != 0 {
		b.Error("should have empty content")
	}

	var keys []string
	var values [][]byte
	for i := 0; i < 10; i++ {
		k := randBytes(30)
		v := randBytes(8)
		keys = append(keys, string(k))
		values = append(values, v)
	}
	one.update(keys, values)

	for i := 0; i < b.N; i++ {
		one.load()
	}

	if len(one.content) != 10 {
		b.Error("should have 10 content")
	}
}

func BenchmarkUpdate(b *testing.B) {
	benchInfo := "bench_info_update"
	defer os.Remove(benchInfo)
	one, err := createInfo(benchInfo)
	if err != nil {
		b.Fatal(err)
	}
	if len(one.content) != 0 {
		b.Error("should have empty content")
	}

	var keys []string
	var values [][]byte
	for i := 0; i < 10; i++ {
		k := randBytes(30)
		v := randBytes(8)
		keys = append(keys, string(k))
		values = append(values, v)
	}

	for i := 0; i < b.N; i++ {
		one.update(keys, values)
	}

	if len(one.content) != 10 {
		b.Error("should have 10 content")
	}
}
