package bytesbufferpool

import (
	"bytes"
	"sync"
)

var defaultPool = &sync.Pool{
	New: func() any {
		return new(bytes.Buffer)
	},
}

func Get() *bytes.Buffer {
	return defaultPool.Get().(*bytes.Buffer)
}

func Put(x *bytes.Buffer) {
	x.Reset()
	defaultPool.Put(x)
}
