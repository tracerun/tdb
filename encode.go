package tdb

import (
	"fmt"
	"path/filepath"
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
	folder := filepath.Join(strconv.Itoa(f.year()), strconv.Itoa(f.month()))
	fileName := strconv.Itoa(f.day())
	if !f.isAM() {
		fileName = fmt.Sprintf("%sp", fileName)
	}
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
