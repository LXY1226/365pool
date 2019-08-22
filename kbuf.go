package main

import (
	"io"
	"time"
)

const buffSize = 1 << 20 //1MiB

type kazybuf struct {
	rd  io.ReadCloser
	buf []byte
	// Cursors always be potion to be read or write
	wCursor  int
	rCursor  int
	isClosed bool
}

func (k kazybuf) Read(p []byte) (n int, err error) {
	if k.isClosed {
		return 0, io.EOF
	}
	for k.wCursor == k.rCursor {
		time.Sleep(250 * time.Millisecond)
	}
	nextCursor := k.rCursor + len(p)
	if nextCursor <= buffSize {
		if nextCursor <= k.wCursor {
			p = k.buf[k.rCursor : k.wCursor-1]
			k.rCursor = k.wCursor
			return k.wCursor - k.rCursor, nil
		} else {
			p = k.buf[k.rCursor : nextCursor-1]
			k.rCursor = nextCursor
			return len(p), nil
		}
	} else {
		if k.wCursor > k.rCursor {
			p = k.buf[k.rCursor : k.wCursor-1]
			k.rCursor = k.wCursor
			return k.wCursor - k.rCursor, nil
		} else {

		}
	}
	return
}

func (k kazybuf) Close() error {
	if k.isClosed {
		return io.EOF
	}
	k.isClosed = true
	k.buf = nil
	return nil
}

func newkazybuf(rd io.ReadCloser) io.ReadCloser {
	var buff kazybuf
	buff.buf = make([]byte, buffSize)
	return io.ReadCloser(buff)
}
