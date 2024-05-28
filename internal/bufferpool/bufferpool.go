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

func Get() *bytes.Buffer {
	b := pool.Get().(*bytes.Buffer)
	return b
}

func Put(b *bytes.Buffer) {
	b.Reset()
	pool.Put(b)
}
