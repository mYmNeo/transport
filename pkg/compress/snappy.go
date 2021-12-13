package compress

import (
	"net/http"

	snappypkg "github.com/golang/snappy"

	"github.com/mYmNeo/transport/pkg/reader"
	"github.com/mYmNeo/transport/pkg/utils"
)

const (
	SnappyEncodingName = "snappy"
)

func init() {
	RegisterCompressionRoundTripperBuilder(SnappyEncodingName, snappyBuilder)
}

type snappyRoundTripper struct {
	base http.RoundTripper
}

func snappyBuilder(rt http.RoundTripper) http.RoundTripper {
	return &snappyRoundTripper{
		base: rt,
	}
}

func (r *snappyRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("Accept-Encoding", utils.AppendContentEncoding(req.Header.Get("Accept-Encoding"), SnappyEncodingName))
	resp, err := r.base.RoundTrip(req)
	if err != nil {
		return nil, err
	}

	matched := utils.MatchContentEncoding(resp.Header.Get("Content-Encoding"), SnappyEncodingName)
	if !matched {
		return resp, nil
	}

	resp.Header.Del("Content-Encoding")

	resp.Body = reader.NewCompressReader(resp.Body, snappypkg.NewReader(resp.Body))

	return resp, nil
}
