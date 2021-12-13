package random

import (
	"bytes"
	"io"
	"math/rand"
	"time"
)

var (
	letters = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	src     = rand.NewSource(time.Now().UnixNano())
)

const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

type randomData struct {
	data *bytes.Buffer
}

type ResetStream interface {
	io.Writer
	io.Closer
	Reset(w io.Writer)
}

func NewRandomGenerator(size int64, w ResetStream) (*randomData, error) {
	r := &randomData{
		data: bytes.NewBuffer(nil),
	}

	b := make([]byte, size)

	for i, cache, remain := size-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letters) {
			b[i] = letters[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	w.Reset(r.data)
	defer w.Close()
	_, err := io.CopyN(w, bytes.NewReader(b), size)
	if err != nil {
		return nil, err
	}

	return r, nil
}

func (g *randomData) GetData() []byte {
	return g.data.Bytes()
}
