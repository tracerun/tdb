package tdb

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInfoMethods(t *testing.T) {
	// create a bad folder
	folderPath := "bad_folder"

	if err := createFolder(folderPath); err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(folderPath)

	// creation on folder will fail
	_, err := createInfo(folderPath)
	assert.Equal(t, ErrPathNotFile, err, "error should pop while creating info on a folder")

	demoInfo := "demo_info"
	defer os.Remove(demoInfo)
	one, err := createInfo(demoInfo)
	assert.NoError(t, err, "should have no error while creating info from non-existed file")
	assert.Len(t, one.content, 0, "should have empty content")

	key := "key"
	value := []byte("value")
	err = one.updateInfo([]string{key}, [][]byte{value})
	assert.NoError(t, err, "should have no error while updating the content")
	assert.Len(t, one.content, 1, "should have one content")

	// get value
	value1 := one.getInfoValue(key)
	assert.EqualValues(t, value, value1, "get value wrong")

	// get a not existed value
	value2 := one.getInfoValue("not exist")
	assert.Nil(t, value2, "should have nil value")

	// get keys
	keys := []string{key}
	assert.EqualValues(t, keys, one.getInfoKeys(), "keys not correct.")

	// get all
	allKeys := []string{key}
	allValues := [][]byte{value}
	keysResult, valuesResult := one.getAllInfo()
	assert.EqualValues(t, allKeys, keysResult, "keys not correct.")
	assert.EqualValues(t, allValues, valuesResult, "values not correct.")

	// k, v length not equal
	err = one.updateInfo([]string{key}, [][]byte{value, value})
	assert.Error(t, err, "should have error when k, v length not equal")

	// empty k, v
	err = one.updateInfo([]string{}, [][]byte{})
	assert.NoError(t, err, "should have no error while updating nothing")

	// load an existing content
	two, err2 := createInfo(demoInfo)
	assert.NoError(t, err2, "should have no error while creating info from an existed file")
	assert.Len(t, two.content, 1, "should have empty content")
	assert.EqualValues(t, one.content[key], two.content[key], "the values should be the same")
}

func BenchmarkLoad10(b *testing.B) {
	benchmarkLoad(b, 10)
}

func BenchmarkUpdate10(b *testing.B) {
	benchmarkUpdate(b, 10)
}

func benchmarkLoad(b *testing.B, count int) {
	benchInfo := "bench_info_load"
	defer os.Remove(benchInfo)
	one, err := createInfo(benchInfo)
	if err != nil {
		b.Fatal(err)
	}

	var keys []string
	var values [][]byte
	for i := 0; i < count; i++ {
		k := randBytes(30)
		v := randBytes(8)
		keys = append(keys, string(k))
		values = append(values, v)
	}
	one.updateInfo(keys, values)

	for i := 0; i < b.N; i++ {
		one.loadInfo()
	}

	if len(one.content) != count {
		b.Errorf("should have %d content", count)
	}
}

func benchmarkUpdate(b *testing.B, count int) {
	benchInfo := "bench_info_update"
	defer os.Remove(benchInfo)
	one, err := createInfo(benchInfo)
	if err != nil {
		b.Fatal(err)
	}

	var keys []string
	var values [][]byte
	for i := 0; i < count; i++ {
		k := randBytes(30)
		v := randBytes(8)
		keys = append(keys, string(k))
		values = append(values, v)
	}

	for i := 0; i < b.N; i++ {
		one.updateInfo(keys, values)
	}

	if len(one.content) != count {
		b.Errorf("should have %d content", count)
	}
}
