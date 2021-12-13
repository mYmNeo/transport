package transport

import (
	"net/http"

	"github.com/mYmNeo/transport/pkg/compress"
)

// NewCompressionRoundTripper returns a http.RoundTripper with supported compression method
func NewCompressionRoundTripper(base http.RoundTripper) http.RoundTripper {
	newRt := base
	for _, builder := range compress.GetBuilders() {
		newRt = builder(newRt)
	}

	return newRt
}
