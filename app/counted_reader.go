package app

import "io"

type countedReader struct {
	r        io.Reader
	callback func(int64)
	sum      int64
}

func (cr *countedReader) Read(p []byte) (int, error) {
	cr.sum += int64(len(p))
	cr.callback(cr.sum)
	return cr.r.Read(p)
}

func NewCountedReader(r io.Reader, callback func(int64)) io.Reader {
	return &countedReader{
		r:        r,
		callback: callback,
	}
}
