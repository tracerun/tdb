package tdb

import (
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"golang.org/x/sync/errgroup"
)

func TestFileLockMethods(t *testing.T) {
	targets := []string{"a", "b", "c", "d"}

	// no lock should create
	file := encodeFileFromUnix(0)
	readLock(targets[0], file)()

	assert.Equal(t, 0, len(slotFileLocker.files), "not loceker should be created")

	// run multiple goroutines to test lock
	var g errgroup.Group
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 5000; i++ {
		num := r.Intn(4)
		read := num%2 == 0
		t := time.Duration(r.Intn(30)) * time.Microsecond
		g.Go(func() error {
			singleAction(targets[num], read, t)
			return nil
		})
	}
	err := g.Wait()
	assert.NoError(t, err, "parallel jobs should have no error")

	assert.Equal(t, 4, len(slotFileLocker.files), "4 lockers should be created")
}

func singleAction(target string, read bool, t time.Duration) {
	file := currentFileEncode()
	if read {
		defer readLock(target, file)()
		time.Sleep(t)
	} else {
		defer writeLock(target, file)()
		time.Sleep(t)
	}
}

func BenchmarkFileLocker(b *testing.B) {
	targets := []string{"a", "b", "c", "d"}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i < b.N; i++ {
		num := r.Intn(4)
		read := num%2 == 0
		singleAction(targets[num], read, 0)
	}
}
