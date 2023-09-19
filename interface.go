package gma

import "io"

type GMA interface {
	Name() string
	Author() string
	Description() string
	Type() string
	Tags() []string

	Entries() []Entry
	Open(path string) (io.ReadSeekCloser, error)
	Close() error
}

type Entry interface {
	Filename() string
	Basename() string
	Path() string
	Length() uint32
	Open() io.ReadSeekCloser
}
