package compress

import (
	"net/http"

	"github.com/mYmNeo/transport/pkg/reader"
	"github.com/mYmNeo/transport/pkg/utils"

	lz4pkg "github.com/pierrec/lz4/v4"
)

const (
	LZ4EncodingName = "lz4"
)

func init() {
	RegisterCompressionRoundTripperBuilder(LZ4EncodingName, lz4Builder)
}

type lz4RoundTripper struct {
	base http.RoundTripper
}

func lz4Builder(rt http.RoundTripper) http.RoundTripper {
	return &lz4RoundTripper{
		base: rt,
	}
}

func (r *lz4RoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("Accept-Encoding", utils.AppendContentEncoding(req.Header.Get("Accept-Encoding"), LZ4EncodingName))
	resp, err := r.base.RoundTrip(req)
	if err != nil {
		return nil, err
	}

	matched := utils.MatchContentEncoding(resp.Header.Get("Content-Encoding"), LZ4EncodingName)
	if !matched {
		return resp, nil
	}

	resp.Header.Del("Content-Encoding")
	resp.Body = reader.NewCompressReader(resp.Body, lz4pkg.NewReader(resp.Body))

	return resp, nil
}
