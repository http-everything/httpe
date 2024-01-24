package buffer

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBufferWrite(t *testing.T) {
	buf := Buffer{}

	// Write some bytes
	n, err := buf.Write([]byte("hello"))
	if n != 5 || err != nil {
		t.Errorf("Write returned unexpected values: %d, %v", n, err)
	}

	// Check buffer contents
	assert.Equal(t, "hello", buf.String())
}

func TestBufferConcurrent(t *testing.T) {
	buf := Buffer{}

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, err := buf.Write([]byte("hello"))
			require.NoError(t, err)
		}()
	}
	wg.Wait()

	assert.Equal(t, 500, buf.Len())
}

func TestBufferCollectingDone(t *testing.T) {
	buf := Buffer{}
	buf.CollectingDone()

	n, err := buf.Write([]byte("data"))
	require.NoError(t, err)
	assert.Equal(t, 4, n)
}
