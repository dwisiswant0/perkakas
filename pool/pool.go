package pool

import (
	"bytes"
	"strings"
	"sync"
)

const (
	TypeBytesBuffer = iota
	TypeStringBuilder
)

var (
	bytesBufferPool = &sync.Pool{
		New: func() interface{} {
			return new(bytes.Buffer)
		},
	}
	stringBuilderPool = &sync.Pool{
		New: func() interface{} {
			return new(strings.Builder)
		},
	}
)

func Get(typePool uint64) interface{} {
	switch typePool {
	case TypeBytesBuffer:
		return bytesBufferPool.Get()
	case TypeStringBuilder:
		return stringBuilderPool.Get()
	default:
		return nil
	}
}

func Put(typePool uint64, inst interface{}) {
	switch typePool {
	case TypeBytesBuffer:
		bytesBufferPool.Put(inst)
	case TypeStringBuilder:
		stringBuilderPool.Put(inst)
	}
}
