package tdb

import (
	"fmt"
	"os"
	"os/user"
	"path"
	"runtime"
	"time"

	"github.com/drkaka/ulid"
)

const (
	metafile = "__metadata__"

	version = uint8(1)

	versionKey  = "version"
	tagKey      = "tag"
	hostKey     = "host"
	archKey     = "arch"
	usernameKey = "username"
	osKey       = "os"
)

func (db *TDB) loadMeta() error {
	var err error
	db.meta, err = createInfo(path.Join(db.path, metafile))
	if err != nil {
		return err
	}
	if len(db.meta.content) == 0 {
		return db.meta.generateMeta()
	}
	return err
}

func (meta *info) generateMeta() error {
	var keys []string
	var values [][]byte
	var err error

	// version
	keys = append(keys, versionKey)
	values = append(values, []byte{byte(version)})

	// tag
	var id ulid.ULID
	if id, err = ulid.NewFromTime(time.Now()); err != nil {
		return err
	}
	keys = append(keys, tagKey)
	values = append(values, id[:])

	// host
	var host string
	if host, err = os.Hostname(); err != nil {
		return err
	}
	keys = append(keys, hostKey)
	values = append(values, []byte(host))

	// arch
	keys = append(keys, archKey)
	values = append(values, []byte(runtime.GOARCH))

	// username
	usr, err := user.Current()
	if err != nil {
		return err
	}
	keys = append(keys, usernameKey)
	values = append(values, []byte(usr.Username))

	// os
	keys = append(keys, osKey)
	values = append(values, []byte(runtime.GOOS))

	return meta.updateInfo(keys, values)
}

func (db *TDB) getMeta(k string) ([]byte, error) {
	if db.meta == nil {
		return nil, ErrNotExist
	}

	b := db.meta.content[k]
	if b == nil {
		return nil, ErrNotExist
	}
	return b, nil
}

// Version for tdb.
func (db *TDB) Version() (uint8, error) {
	b, err := db.getMeta(versionKey)
	if err != nil {
		return 0, err
	}

	if len(b) != 1 {
		return 0, fmt.Errorf("wrong version in metadata, byte length should be 1, but is %d", len(b))
	}

	return uint8(b[0]), nil
}

// Tag for database
func (db *TDB) Tag() (string, error) {
	b, err := db.getMeta(tagKey)
	if err != nil {
		return "", err
	}

	var uid ulid.ULID
	if l := len(b); l != 16 {
		return "", fmt.Errorf("wrong tag in metadata, byte length should be 16, but is %d (%s)", l, string(b))
	}
	copy(uid[:], b[:16])
	return uid.String(), nil
}

// CreateAt unixtime
func (db *TDB) CreateAt() (uint32, error) {
	b, err := db.getMeta(tagKey)
	if err != nil {
		return 0, err
	}

	var uid ulid.ULID
	if l := len(b); l != 16 {
		return 0, fmt.Errorf("wrong tag in metadata, byte length should be 1, but is %d (%s)", l, string(b))
	}
	copy(uid[:], b[:16])

	return uint32(uid.Time() / 1000), nil
}

// get string meta info
func (db *TDB) stringMeta(k string) (string, error) {
	b, err := db.getMeta(k)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// Host name
func (db *TDB) Host() (string, error) {
	return db.stringMeta(hostKey)
}

// Arch info
func (db *TDB) Arch() (string, error) {
	return db.stringMeta(archKey)
}

// Username info
func (db *TDB) Username() (string, error) {
	return db.stringMeta(usernameKey)
}

// OS info
func (db *TDB) OS() (string, error) {
	return db.stringMeta(osKey)
}
