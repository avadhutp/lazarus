package geddit

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var testJSON = `
	{
		"kind" : "Listing",
		"data" : {
			"after" : "test-after",
			"children" : [
				{
					"kind" : "t3",
					"data" : {
						"domain" 			: "youtube.com",
						"url"	 			: "youtube.com/url/test",
						"title"	 			: "Test song title",
						"link_flair_text" 	: "test-genre",
						"id"				: "12345"
					}
				}
			]
		}
	}
`

func TestGet(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, testJSON)
	}))
	defer ts.Close()

	actual := Get(ts.URL)

	assert.Equal(t, "youtube.com", actual["12345"].Data.Domain)
	assert.Equal(t, "youtube.com/url/test", actual["12345"].Data.URL)
	assert.Equal(t, "test-genre", actual["12345"].Data.Genre)
	assert.Equal(t, "Test song title", actual["12345"].Data.Title)
	assert.Equal(t, "12345", actual["12345"].Data.ID)
}

func TestGetHandleServerErrors(t *testing.T) {
	tests := []struct {
		status   int
		response string
		msg      string
	}{
		{http.StatusInternalServerError, testJSON, "Non-200 response code should yield no response"},
		{http.StatusOK, strings.TrimSuffix("}", testJSON), "Malformed JSON should yield no response"},
	}

	for _, test := range tests {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(test.status)
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprintln(w, test.response)
		}))
		defer ts.Close()

		actual := Get(ts.URL)
		assert.Empty(t, actual, test.msg)
	}
}

func TestGetHandlErrors(t *testing.T) {
	actual := Get(httptest.DefaultRemoteAddr)

	assert.Empty(t, actual)
}
