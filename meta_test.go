package tdb

import (
	"os"
	"os/user"
	"runtime"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMetadata(t *testing.T) {
	metaTestFolder := "meta_test_folder"
	db, err := Open(metaTestFolder)
	assert.NoError(t, err, "should have no error opening TDB")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(metaTestFolder)

	now := uint32(time.Now().Unix())
	// test version
	v, err := db.Version()
	assert.NoError(t, err, "should have no error to get db version")
	assert.Equal(t, version, v, "version wrong")

	// test tag
	tag, err := db.Tag()
	assert.NoError(t, err, "should have no error to get db tag")
	assert.Len(t, tag, 26, "ulid tag length should be 26")

	// test creation time
	createAt, err := db.CreateAt()
	assert.NoError(t, err, "should have no error to get db creation")
	assert.InDelta(t, now, createAt, 1, "creation is wrong")

	// test host
	host, err := db.Host()
	assert.NoError(t, err, "should have no error to get db host")

	if thisHost, err := os.Hostname(); err != nil {
		t.Error(err)
	} else {
		assert.EqualValues(t, host, thisHost, "host not correct")
	}

	// test user name
	username, err := db.Username()
	if usr, err := user.Current(); err != nil {
		t.Error(err)
	} else {
		assert.EqualValues(t, usr.Username, username, "username not correct")
	}

	// test arch
	arch, err := db.Arch()
	assert.NoError(t, err, "should have no error to get db arch")
	assert.EqualValues(t, runtime.GOARCH, arch, "arch not correct")

	// test OS
	os, err := db.OS()
	assert.NoError(t, err, "should have no error to get db os")
	assert.EqualValues(t, runtime.GOOS, os, "os not correct")

	// test zone offset
	_, thisOffset := time.Now().Zone()
	offset, err := db.ZoneOffset()
	assert.NoError(t, err, "should have no error to get zone offset")
	assert.Equal(t, int32(thisOffset), offset, "offset value is wrong")
}
