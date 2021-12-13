package compress

import (
	"net/http"
)

var (
	builders = make(map[string]BuildCompressionRoundTripperFunc)
)

type BuildCompressionRoundTripperFunc func(rt http.RoundTripper) http.RoundTripper

// RegisterCompressionRoundTripperBuilder register a compression round tripper
func RegisterCompressionRoundTripperBuilder(name string, builder BuildCompressionRoundTripperFunc) {
	if _, ok := builders[name]; !ok {
		builders[name] = builder
	}
}

// GetBuilders returns a map of compression roundtripper builders
func GetBuilders() map[string]BuildCompressionRoundTripperFunc {
	return builders
}
