package bufferpool

import (
	"bytes"
	"sync"
)

var pool = sync.Pool{
	New: func() any {
		return new(bytes.Buffer)
	},
}

func GetBuffer() *bytes.Buffer {
	b := pool.Get().(*bytes.Buffer)
	return b
}

func PutBuffer(b *bytes.Buffer) {
	b.Reset()
	pool.Put(b)
}
