package tdb

import "math/rand"

const (
	letterBytes = "0123456789abcdefghijklmnopqrstuvwxyz"
	letterLen   = 36
)

func randBytes(n int) []byte {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(letterLen)]
	}
	return b
}
