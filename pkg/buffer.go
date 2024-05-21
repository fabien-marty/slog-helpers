package slogh

import (
	"bytes"
	"sync"
)

var bufPool = sync.Pool{
	New: func() any {
		return new(bytes.Buffer)
	},
}

func getBuffer() *bytes.Buffer {
	b := bufPool.Get().(*bytes.Buffer)
	return b
}

func putBuffer(b *bytes.Buffer) {
	b.Reset()
	bufPool.Put(b)
}
