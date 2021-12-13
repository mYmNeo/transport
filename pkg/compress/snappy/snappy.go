package snappy

import (
	"net/http"

	snappypkg "github.com/golang/snappy"

	"github.com/mYmNeo/transport/pkg/reader"
	"github.com/mYmNeo/transport/pkg/rt"
	"github.com/mYmNeo/transport/pkg/utils"
)

const (
	EncodingName = "snappy"
)

func init() {
	rt.RegisterCompressionRoundTripperBuilder(EncodingName, builder)
}

type snappyRoundTripper struct {
	base http.RoundTripper
}

func builder(rt http.RoundTripper) http.RoundTripper {
	return &snappyRoundTripper{
		base: rt,
	}
}

func (r *snappyRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("Accept-Encoding", utils.AppendContentEncoding(req.Header.Get("Accept-Encoding"), EncodingName))
	resp, err := r.base.RoundTrip(req)
	if err != nil {
		return nil, err
	}

	matched := utils.MatchContentEncoding(resp.Header.Get("Content-Encoding"), EncodingName)
	if !matched {
		return resp, nil
	}

	resp.Header.Del("Content-Encoding")

	resp.Body = reader.NewCompressReader(resp.Body, snappypkg.NewReader(resp.Body))

	return resp, nil
}
