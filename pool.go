package hbuffer

import (
	"sync"
)

type BufferPool struct {
	pool sync.Pool
}

func NewBufferPool() *BufferPool {
	return &BufferPool{pool: sync.Pool{New: func() interface{} {
		return NewBuffer()
	}}}
}

func (b *BufferPool) Get() *Buffer {
	return b.pool.Get().(*Buffer)
}

func (b *BufferPool) Put(buf *Buffer) {
	if buf == nil {
		return
	}
	buf.Reset()
	b.pool.Put(buf)
}
