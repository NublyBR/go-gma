package gma

import (
	"io"
	"path/filepath"
)

type entry struct {
	parent *gma

	name string
	size uint32
	offs int64
}

func (e *entry) Filename() string {
	return e.name
}

func (e *entry) Basename() string {
	return filepath.Base(e.name)
}

func (e *entry) Path() string {
	return filepath.Dir(e.name)
}

func (e *entry) Length() uint32 {
	return e.size
}

func (e *entry) Open() io.ReadSeekCloser {
	return &entryReader{
		fs:       e.parent.stream,
		closed:   false,
		offset:   e.parent.start + e.offs,
		size:     int64(e.size),
		position: 0,
	}
}
