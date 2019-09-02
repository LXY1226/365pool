package main

import (
	"bytes"
	"io"
	"time"
)

const buffSize = 1 << 20 //1MiB

/*
1.	----------------------------------
	      ^w              ^r


*/
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
		nextCursor -= buffSize
		if k.wCursor > k.rCursor {
			p = k.buf[k.rCursor : k.wCursor-1]
			k.rCursor = k.wCursor
			return k.wCursor - k.rCursor, nil
		} else {
			if nextCursor <= k.wCursor {
				var buf bytes.Buffer
				buf.Write(k.buf[k.rCursor:])
				buf.Write(k.buf[0:nextCursor])
				k.rCursor = nextCursor
				p = buf.Bytes()
				return len(p), nil
			}
		}
	}
	return
}

func (k kazybuf) readloop() {
	for {
		iBuff := make([]byte, buffSize/4)
		n, err := k.rd.Read(iBuff)
		if err != nil {
			if err == io.EOF {
				_ = k.Close()
			}
		}
		// Wait For Usable Buffer
		for k.wCursor == k.rCursor-1 {
			time.Sleep(250 * time.Millisecond)
		}
		i := 0
		spaceMore := k.rCursor - k.wCursor - 1
		if k.wCursor < k.rCursor {
			i = k.rCursor - 1 - k.wCursor
			if i > n {
				i = n
			}

			spaceMore += buffSize
		}

	}
}

func (k kazybuf) Close() error {
	if k.isClosed {
		return io.EOF
	}
	k.isClosed = true
	// Giving up unused data
	k.buf = nil
	return nil
}

func (k kazybuf) reInit() {

}

func newkazybuf(rd io.ReadCloser) io.ReadCloser {
	var buff kazybuf
	buff.buf = make([]byte, buffSize)
	buff.rd = rd
	go buff.readloop()
	return buff
}
