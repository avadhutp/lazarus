package ui

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestExpandYoutubeURL(t *testing.T) {
	tests := []struct {
		in       string
		expected string
		msg      string
	}{
		{
			in:       "https://youtu.be/zKVbJuhcze8",
			expected: "http://www.youtube.com/watch?v=zKVbJuhcze8",
			msg:      "Shortened youtu.be url should be expanded correctly",
		},
		{
			in:       "http://[fe80::%31%25en0]:8080/",
			expected: "http://[fe80::%31%25en0]:8080/",
			msg:      "Invalid url fails url.Parse and therefore should be returned as is",
		},
		{
			in:       "http://www.soundcloud.com/test",
			expected: "http://www.soundcloud.com/test",
			msg:      "Not a youtube URL, so should be returned as is",
		},
	}

	for _, test := range tests {
		actual := expandYoutubeURL(test.in)
		assert.Equal(t, test.expected, actual, test.msg)
	}
}
