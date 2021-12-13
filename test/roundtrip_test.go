package test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	snappyw "github.com/golang/snappy"
	lz4w "github.com/pierrec/lz4/v4"

	"github.com/mYmNeo/transport/pkg/compress/lz4"
	"github.com/mYmNeo/transport/pkg/compress/snappy"
	"github.com/mYmNeo/transport/pkg/rt"
	"github.com/mYmNeo/transport/test/random"
)

func TestLZ4RoundTripper(t *testing.T) {
	testSize := []int64{1 << 10, 1 << 20, 1 << 24}

	for _, size := range testSize {
		generator, err := random.NewRandomGenerator(size, lz4w.NewWriter(nil))
		if err != nil {
			t.Errorf("can't create generator")
			continue
		}

		func() {
			s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				w.Header().Set("Content-Encoding", lz4.EncodingName)
				w.WriteHeader(http.StatusOK)

				w.Write(generator.GetData())
			}))

			defer s.Close()
			client := s.Client()

			newRt := rt.NewCompressionRoundTripper(client.Transport)
			client.Transport = newRt

			resp, err := client.Get(s.URL)
			if err != nil {
				t.Errorf("can't perform request, %s", err.Error())
				return
			}

			if resp.StatusCode != http.StatusOK {
				t.Errorf("status is not 200, got %d", resp.StatusCode)
				return
			}

			if v := resp.Header.Get("Content-Encoding"); v != "" {
				t.Errorf("should not have Content-Encoding %s", v)
				return
			}

			defer resp.Body.Close()
			n, err := io.Copy(io.Discard, resp.Body)
			if err != nil {
				t.Errorf("read failed, %s", err.Error())
				return
			}

			if n != size {
				t.Errorf("expected %d bytes, got %d bytes", size, n)
				return
			}
		}()
	}
}

func TestSnappyRoundTripper(t *testing.T) {
	testSize := []int64{1 << 10, 1 << 20, 1 << 24}

	for _, size := range testSize {
		generator, err := random.NewRandomGenerator(size, snappyw.NewBufferedWriter(nil))
		if err != nil {
			t.Errorf("can't create generator")
			continue
		}

		func() {
			s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				w.Header().Set("Content-Encoding", snappy.EncodingName)
				w.WriteHeader(http.StatusOK)

				w.Write(generator.GetData())
			}))

			defer s.Close()
			client := s.Client()

			client.Transport = rt.NewCompressionRoundTripper(client.Transport)

			resp, err := client.Get(s.URL)
			if err != nil {
				t.Errorf("can't perform request, %s", err.Error())
				return
			}

			if resp.StatusCode != http.StatusOK {
				t.Errorf("status is not 200, got %d", resp.StatusCode)
				return
			}

			if v := resp.Header.Get("Content-Encoding"); v != "" {
				t.Errorf("should not have Content-Encoding %s", v)
				return
			}

			defer resp.Body.Close()
			n, err := io.Copy(io.Discard, resp.Body)
			if err != nil {
				t.Errorf("read failed, %s", err.Error())
				return
			}

			if n != size {
				t.Errorf("expected %d bytes, got %d bytes", size, n)
				return
			}
		}()
	}
}
