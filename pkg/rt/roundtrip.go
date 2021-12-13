package rt

import (
	"net/http"
)

var (
	builders map[string]BuildCompressionRoundTripperFunc
)

type BuildCompressionRoundTripperFunc func(rt http.RoundTripper) http.RoundTripper

func init() {
	builders = make(map[string]BuildCompressionRoundTripperFunc)
}

// RegisterCompressionRoundTripperBuilder register a compression round tripper
func RegisterCompressionRoundTripperBuilder(name string, builder BuildCompressionRoundTripperFunc) {
	if _, ok := builders[name]; !ok {
		builders[name] = builder
	}
}

// NewCompressionRoundTripper returns a http.RoundTripper with supported compression method
func NewCompressionRoundTripper(base http.RoundTripper) http.RoundTripper {
	newRt := base
	for _, builder := range builders {
		newRt = builder(newRt)
	}

	return newRt
}
