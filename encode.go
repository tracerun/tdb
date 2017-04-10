package tdb

import (
	"math"
	"strconv"
	"time"
)

type fileEncode uint32

func encodeFileFromUnix(unixtime uint32) fileEncode {
	t := time.Unix(int64(unixtime), 0)
	year, month, day := t.Date()
	hour := t.Hour()
	return encodeFile(year, int(month), day, hour)
}

func encodeFile(year, month, day, hour int) fileEncode {
	var encoded uint32
	encoded += uint32(year) * 1e5
	encoded += uint32(month) * 1e3
	encoded += uint32(day) * 10
	if hour >= 12 {
		encoded += 5
	}
	return fileEncode(encoded)
}

func (f fileEncode) year() int {
	return int(uint32(f) / 1e5)
}

func (f fileEncode) month() int {
	return int(uint32(f)/1e3) % 100
}

func (f fileEncode) day() int {
	return int(uint32(f)/10) % 100
}

func (f fileEncode) isAM() bool {
	return f%10 == 0
}

// get encoded folder and file name
func (f fileEncode) path() (string, string) {
	folder := strconv.Itoa(int(uint32(f) / 1e3))
	fileName := strconv.Itoa(int(uint32(f) % 1e3))
	return folder, fileName
}

// get original unix time
func (f fileEncode) origin() uint32 {
	fileOrigin := time.Date(f.year(), time.Month(f.month()), f.day(), 0, 0, 0, 0, time.Local).Unix()
	if !f.isAM() {
		fileOrigin = fileOrigin + 43200
	}
	return uint32(fileOrigin)
}

type fileRange struct {
	start fileEncode
	end   fileEncode
}

// create a new fileRange instance
// startUnix == 0 means from very beginning; endUnix == 0 means to very end.
func newFileRange(startUnix uint32, endUnix uint32) (*fileRange, error) {
	if startUnix != 0 && endUnix != 0 && startUnix > endUnix {
		return nil, ErrRange
	}

	one := &fileRange{
		start: encodeFileFromUnix(startUnix),
		end:   encodeFileFromUnix(endUnix),
	}
	// set to 0 if startUnix == 0
	if startUnix == 0 {
		one.start = fileEncode(0)
	}
	// set to max uint32 if endUnix == 0
	if endUnix == 0 {
		one.end = fileEncode(math.MaxUint32)
	}
	return one, nil
}

// folder should be the year-month folder like "201703"
func (f *fileRange) folderInRange(folder string) (bool, error) {
	v, err := strconv.Atoi(folder)
	if err != nil {
		return false, err
	}

	t := uint32(v)
	if t >= uint32(f.start/1000) && t <= uint32(f.end/1000) {
		return true, nil
	}
	return false, nil
}
