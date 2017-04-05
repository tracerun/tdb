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

	if stat, err := os.Stat(p); err != nil {
		if os.IsNotExist(err) {
			return one, nil
		}
		return one, err
	} else if stat.IsDir() {
		return one, ErrInfoFilePath
	}

	err := one.loadContent()

	return one, err
}

func (one *info) loadContent() error {
	one.fileLock.RLock()
	b, _ := ioutil.ReadFile(one.path)
	one.fileLock.RUnlock()

	one.contentLock.Lock()
	defer one.contentLock.Unlock()

	var pbInfo Info
	if err := proto.Unmarshal(b, &pbInfo); err != nil {
		return err
	}
	one.content = pbInfo.Fields

	return nil
}

func (one *info) setContent(k []string, v [][]byte) error {
	if len(k) != len(v) {
		return errors.New("k, v length not equal")
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
