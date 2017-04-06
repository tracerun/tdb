package tdb

import "math/rand"

const (
	letterBytes = "123456789abcdefghijklmnopqrstuvwxyz0"
	letterLen   = 36
)

func randBytes(n int) []byte {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(letterLen)]
	}
	return b
}

func addSlot(index *info, target string, start, howlong uint32) {

}
