package buffer

import (
	"bytes"
	"sync"
	"sync/atomic"
)

// Buffer is a goroutine safe bytes.Buffer
type Buffer struct {
	buffer         bytes.Buffer
	mutex          sync.Mutex
	stopCollecting atomic.Bool
}

func (b *Buffer) Write(p []byte) (n int, err error) {
	if b.stopCollecting.Load() {
		return len(p), nil
	}
	b.mutex.Lock()
	defer b.mutex.Unlock()
	return b.buffer.Write(p)
}

func (b *Buffer) CollectingDone() {
	b.stopCollecting.Store(true)
}

func (b *Buffer) String() string {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	return b.buffer.String()
}

func (b *Buffer) Len() int {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	return b.buffer.Len()
}
