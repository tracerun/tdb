package tdb

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRandBytes(t *testing.T) {
	assert.Len(t, randBytes(10), 10, "bytes length should be 10")
}

func BenchmarkRandBytes(b *testing.B) {
	for i := 0; i < b.N; i++ {
		randBytes(10)
	}
}
