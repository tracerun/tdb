package tdb

import (
	"errors"
	"io/ioutil"
	"os"
	"sync"

	"github.com/golang/protobuf/proto"
)

type info struct {
	path        string
	content     map[string][]byte
	contentLock *sync.RWMutex
	fileLock    *sync.RWMutex
}

func createInfo(p string) (*info, error) {
	one := &info{path: p}
	one.content = make(map[string][]byte)
	one.contentLock = new(sync.RWMutex)
	one.fileLock = new(sync.RWMutex)

	if stat, err := os.Stat(p); err != nil {
		if os.IsNotExist(err) {
			return one, nil
		}
		return one, err
	} else if stat.IsDir() {
		return one, ErrPathNotFile
	}

	err := one.loadInfo()

	return one, err
}

// loadInfo information from file.
func (one *info) loadInfo() error {
	one.contentLock.Lock()
	defer one.contentLock.Unlock()

	one.fileLock.RLock()
	b, err := ioutil.ReadFile(one.path)
	one.fileLock.RUnlock()
	if err != nil {
		return err
	}

	var pbInfo Info
	if err := proto.Unmarshal(b, &pbInfo); err != nil {
		return err
	}

	one.content = pbInfo.Fields

	return nil
}

func (one *info) getInfoValue(k string) []byte {
	one.contentLock.RLock()
	defer one.contentLock.RUnlock()
	return one.content[k]
}

func (one *info) getInfoKeys() []string {
	one.contentLock.RLock()
	defer one.contentLock.RUnlock()

	var keys []string
	for k := range one.content {
		keys = append(keys, k)
	}
	return keys
}

func (one *info) getAllInfo() ([]string, [][]byte) {
	one.contentLock.RLock()
	defer one.contentLock.RUnlock()

	var keys []string
	var values [][]byte
	for k, v := range one.content {
		keys = append(keys, k)
		values = append(values, v)
	}
	return keys, values
}

// update the given keys, values content and write to file
func (one *info) updateInfo(k []string, v [][]byte) error {
	if len(k) != len(v) {
		return errors.New("k, v length not equal")
	}

	if len(k) == 0 {
		return nil
	}

	one.contentLock.Lock()
	for i := 0; i < len(k); i++ {
		one.content[k[i]] = v[i]
	}
	var pbInfo Info
	pbInfo.Fields = one.content
	one.contentLock.Unlock()

	b, err := proto.Marshal(&pbInfo)
	if err != nil {
		return err
	}

	one.fileLock.Lock()
	defer one.fileLock.Unlock()

	return ioutil.WriteFile(one.path, b, 0644)
}
