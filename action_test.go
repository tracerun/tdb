package tdb

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestActions(t *testing.T) {
	folder := "test_actions"
	db, err := Open(folder)
	if !assert.NoError(t, err, "should have no error") {
		t.Fatal(err)
	}
	defer os.RemoveAll(folder)

	testNormalCloseAction(t, db)
	testSingleCloseAction(t, db)
	testCloseEarlier(t, db)
	testCloseExpiredAction(t, db)
	testLaterActiveAction(t, db)

	// test load actions
	db2, err := Open(folder)
	if !assert.NoError(t, err, "should have no error") {
		t.Fatal(err)
	}

	targets, allStarts, allLasts, err := db2.GetActions()
	assert.NoError(t, err, "get all actions wrong")
	assert.Len(t, targets, 1, "targets count wrong")
	assert.Len(t, allStarts, 1, "starts count wrong")
	assert.Len(t, allLasts, 1, "lasts count wrong")
}

func testNormalCloseAction(t *testing.T, db *TDB) {
	target := string(randBytes(6))
	now := uint32(time.Now().Unix())

	// close after 2 seconds
	db.AddAction(target, true, now)
	db.AddAction(target, false, now+2)

	starts, howlongs, err := db.GetAllSlots(target)
	assert.NoError(t, err, "get slots wrong")
	assert.Len(t, starts, 1, "starts count wrong")
	assert.Len(t, howlongs, 1, "howlong count wrong")

	assert.Equal(t, now, starts[0][0], "start wrong")
	assert.Equal(t, uint32(2), howlongs[0][0], "last wrong")
}

func testSingleCloseAction(t *testing.T, db *TDB) {
	target := string(randBytes(6))
	now := uint32(time.Now().Unix())

	// single close action
	db.AddAction(target, false, now)

	starts, howlongs, err := db.GetAllSlots(target)
	assert.NoError(t, err, "get slots wrong")
	assert.Len(t, starts, 1, "starts count wrong")
	assert.Len(t, howlongs, 1, "howlong count wrong")

	assert.Equal(t, now, starts[0][0], "start wrong")
	assert.Equal(t, uint32(1), howlongs[0][0], "last wrong")
}

func testCloseEarlier(t *testing.T, db *TDB) {
	target := string(randBytes(6))
	now := uint32(time.Now().Unix())

	// close earlier
	db.AddAction(target, true, now)
	db.AddAction(target, true, now+3)
	db.AddAction(target, false, now+2)

	starts, howlongs, err := db.GetAllSlots(target)
	assert.NoError(t, err, "get slots wrong")
	assert.Len(t, starts, 1, "starts count wrong")
	assert.Len(t, howlongs, 1, "howlong count wrong")

	assert.Equal(t, now, starts[0][0], "start wrong")
	assert.Equal(t, uint32(3), howlongs[0][0], "last wrong")
}

func testCloseExpiredAction(t *testing.T, db *TDB) {
	target := string(randBytes(6))
	now := uint32(time.Now().Unix())

	// close expired
	db.AddAction(target, true, now)
	db.AddAction(target, false, now+17)

	starts, howlongs, err := db.GetAllSlots(target)
	assert.NoError(t, err, "get slots wrong")
	assert.Len(t, starts[0], 2, "starts count wrong")
	assert.Len(t, howlongs[0], 2, "howlong count wrong")

	assert.Equal(t, now, starts[0][0], "start wrong")
	assert.Equal(t, uint32(1), howlongs[0][0], "last wrong")

	assert.Equal(t, now+17, starts[0][1], "start wrong")
	assert.Equal(t, uint32(1), howlongs[0][1], "last wrong")

	targets, allStarts, allLasts, err := db.GetActions()
	assert.NoError(t, err, "get all actions wrong")
	assert.Len(t, targets, 0, "targets count wrong")
	assert.Len(t, allStarts, 0, "starts count wrong")
	assert.Len(t, allLasts, 0, "lasts count wrong")
}

func testLaterActiveAction(t *testing.T, db *TDB) {
	target := string(randBytes(6))
	now := uint32(time.Now().Unix())

	// later active
	db.AddAction(target, true, now)
	db.AddAction(target, true, now+17)

	starts, howlongs, err := db.GetAllSlots(target)
	assert.NoError(t, err, "get slots wrong")
	assert.Len(t, starts[0], 1, "starts count wrong")
	assert.Len(t, howlongs[0], 1, "howlong count wrong")

	assert.Equal(t, now, starts[0][0], "start wrong")
	assert.Equal(t, uint32(1), howlongs[0][0], "last wrong")

	targets, allStarts, allLasts, err := db.GetActions()
	assert.NoError(t, err, "get all actions wrong")
	assert.Len(t, targets, 1, "targets count wrong")
	assert.Len(t, allStarts, 1, "starts count wrong")
	assert.Len(t, allLasts, 1, "lasts count wrong")
}
