package gma

import (
	"io"
	"os"
)

type entryReader struct {
	fs io.ReadSeekCloser

	closed bool

	offset   int64
	size     int64
	position int64
}

func (e *entryReader) Read(p []byte) (int, error) {
	if e.closed {
		return 0, os.ErrClosed
	}
	if e.position >= e.size {
		return 0, io.EOF
	}

	want := int64(len(p))

	if e.position+want > e.size {
		want = e.size - e.position
	}

	e.fs.Seek(e.offset+e.position, io.SeekStart)
	n, err := e.fs.Read(p[:want])
	e.position += int64(n)

	return n, err
}

func (e *entryReader) Seek(offset int64, whence int) (int64, error) {
	if e.closed {
		return 0, os.ErrClosed
	}

	switch whence {
	case io.SeekStart:
		e.position = offset
	case io.SeekCurrent:
		e.position += offset
	case io.SeekEnd:
		e.position = e.size + offset
	}

	if e.position < 0 {
		e.position = 0
	} else if e.position > e.size {
		e.position = e.size
	}

	return e.position, nil
}

func (e *entryReader) Close() error {
	if e.closed {
		return os.ErrClosed
	}
	e.closed = true
	return nil
}
