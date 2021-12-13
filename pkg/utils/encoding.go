package utils

import (
	"strings"
)

func MatchContentEncoding(encoding string, expected string) bool {
	matched := false

	for len(encoding) > 0 {
		var token string
		if next := strings.Index(encoding, ","); next != -1 {
			token = encoding[:next]
			encoding = encoding[next+1:]
		} else {
			token = encoding
			encoding = ""
		}

		if strings.TrimSpace(token) == expected {
			matched = true
			break
		}
	}

	return matched
}

func AppendContentEncoding(oldEncoding, newEncoding string) string {
	return strings.TrimPrefix(strings.TrimSpace(strings.Join([]string{oldEncoding, newEncoding}, ",")), ",")
}
