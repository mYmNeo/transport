package reader

import (
	"io"
)

type compressReader struct {
	body    io.ReadCloser
	wrapper io.Reader
}

func NewCompressReader(respBody io.ReadCloser, cReader io.Reader) *compressReader {
	return &compressReader{
		body:    respBody,
		wrapper: cReader,
	}
}

func (r *compressReader) Read(p []byte) (n int, err error) {
	return r.wrapper.Read(p)
}

func (r *compressReader) Close() error {
	return r.body.Close()
}
