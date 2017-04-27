package tdb

import (
	"os"
	"testing"
	"time"

	"github.com/drkaka/lg"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

func TestActions(t *testing.T) {
	lg.InitLogger(true)

	folder := "test_actions"
	db, err := Open(folder)
	if !assert.NoError(t, err, "should have no error") {
		t.Fatal(err)
	}
	defer os.RemoveAll(folder)

	testFastInsert(t, db)
	testAddAction(t, db)

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

	// after test, should have no item
	testCheckEXP(t, db)

	testReloadEmptyInfo(t, folder)
}

func testFastInsert(t *testing.T, db *TDB) {
	target := string(randBytes(6))
	now := uint32(time.Now().Unix())

	p("test fast insert", zap.String("target", target))

	var g errgroup.Group
	for i := 0; i < 100; i++ {
		g.Go(func() error {
			return db.AddAction(target, now)
		})
	}
	err := g.Wait()
	assert.NoError(t, err, "error when wait")

	starts, howlongs, err := db.GetSlots(target, 0, 0)
	assert.NoError(t, err, "get slots wrong")
	assert.Len(t, starts, 0, "starts count wrong")
	assert.Len(t, howlongs, 0, "howlong count wrong")

	delete(db.action.content, target)
}

func testAddAction(t *testing.T, db *TDB) {
	target := string(randBytes(6))
	now := uint32(time.Now().Unix())

	// add new action
	db.AddAction(target, now)
	// should create a slot
	db.AddAction(target, now+17)

	starts, howlongs, err := db.GetSlots(target, 0, 0)
	assert.NoError(t, err, "get slots wrong")
	assert.Len(t, starts, 1, "starts count wrong")
	assert.Len(t, howlongs, 1, "howlong count wrong")

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

func testCheckEXP(t *testing.T, db *TDB) {
	// delete all actions
	var keys []string
	for k := range db.action.content {
		keys = append(keys, k)
	}
	for i := 0; i < len(keys); i++ {
		delete(db.action.content, keys[i])
	}

	target := string(randBytes(6))
	now := uint32(time.Now().Unix())

	// add a should be expired action
	err := db.AddAction(target, now-actionExp-2)
	assert.NoError(t, err, "error add an action")

	err = db.CheckExpirations()
	assert.NoError(t, err, "error check expiration")

	targets, allStarts, allLasts, err := db.GetActions()
	assert.NoError(t, err, "get all actions wrong")
	assert.Len(t, targets, 0, "targets count wrong")
	assert.Len(t, allStarts, 0, "starts count wrong")
	assert.Len(t, allLasts, 0, "lasts count wrong")
}

func testReloadEmptyInfo(t *testing.T, folder string) {
	db, err := Open(folder)
	if !assert.NoError(t, err, "should have no error") {
		t.Fatal(err)
	}

	target := string(randBytes(6))
	now := uint32(time.Now().Unix())

	err = db.AddAction(target, now)
	assert.NoError(t, err, "error add an action")
}
