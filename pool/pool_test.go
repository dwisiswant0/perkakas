package pool

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetTypePoolByteBuffer(t *testing.T) {
	poolType := uint64(0)
	assert.Equal(t, Get(poolType), bytesBufferPool.Get())
}

func TestGetTypePoolStringBuilder(t *testing.T) {
	poolType := uint64(1)
	assert.Equal(t, Get(poolType), stringBuilderPool.Get())
}

func TestGetTypePoolDefault(t *testing.T) {
	poolType := uint64(2)
	assert.Equal(t, Get(poolType), nil)
}

func TestPutTypePoolBytesBuffer(t *testing.T) {
	poolType := uint64(0)
	buffer := bytesBufferPool.Get()

	Put(poolType, buffer)
}

func TestPutTypePoolStringBuilder(t *testing.T) {
	poolType := uint64(1)
	buffer := stringBuilderPool.Get()

	Put(poolType, buffer)
}
