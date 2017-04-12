package tdb

import "sync"

type fileLocker struct {
	files map[string]*sync.RWMutex
	l     *sync.RWMutex
}

var (
	slotFileLocker *fileLocker
)

func init() {
	slotFileLocker = &fileLocker{
		files: make(map[string]*sync.RWMutex),
		l:     new(sync.RWMutex),
	}
}

// read lock a certain target, if locker not existed, create one.
func (locker *fileLocker) readLock(alias string) {
	locker.l.RLock()
	thisLocker := locker.files[alias]
	locker.l.RUnlock()

	if thisLocker == nil {
		// file locker not existed, create one
		locker.l.Lock()

		thisLocker = locker.files[alias]
		if thisLocker == nil {
			thisLocker = new(sync.RWMutex)
			locker.files[alias] = thisLocker
		}

		locker.l.Unlock()
	}
	thisLocker.RLock()
}

// write lock a certain target, if locker not existed, create one.
func (locker *fileLocker) writeLock(alias string) {
	locker.l.RLock()
	thisLocker := locker.files[alias]
	locker.l.RUnlock()

	if thisLocker == nil {
		// file locker not existed, create one
		locker.l.Lock()

		thisLocker = locker.files[alias]
		if thisLocker == nil {
			thisLocker = new(sync.RWMutex)
			locker.files[alias] = thisLocker
		}

		locker.l.Unlock()
	}
	thisLocker.Lock()
}

// read unlock a certain target, locker should be existed.
func (locker *fileLocker) readUnlock(alias string) {
	locker.l.RLock()
	thisLocker := locker.files[alias]
	locker.l.RUnlock()

	thisLocker.RUnlock()
}

// write unlock a certain target, locker should be existed.
func (locker *fileLocker) writeUnlock(alias string) {
	locker.l.RLock()
	thisLocker := locker.files[alias]
	locker.l.RUnlock()

	thisLocker.Unlock()
}

// read lock the aliased target for the encoded file
// return cleanup func
func readLock(alias string, file fileEncode) func() {
	currentFile := currentFileEncode()
	if uint32(file) != uint32(currentFile) {
		// if file is not the lastest, no need to lock, only read can happen
		return func() {}
	}

	slotFileLocker.readLock(alias)
	return func() {
		slotFileLocker.readUnlock(alias)
	}
}

// write lock the aliased target for the encoded file
// return cleanup func
func writeLock(alias string, file fileEncode) func() {
	currentFile := currentFileEncode()
	if uint32(file) != uint32(currentFile) {
		// if file is not the lastest, no need to lock, only read can happen
		return func() {}
	}

	slotFileLocker.writeLock(alias)
	return func() {
		slotFileLocker.writeUnlock(alias)
	}
}
