package geddit

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// Get Hit the feed URL and get a struct of items ready for display/download/play in lazarus
func Get(rURL string) (after string, music map[string]*Children) {
	r, err := http.Get(rURL)

	if err != nil {
		fmt.Errorf("Unable to get the subreddit %s; error faced: %s", rURL, err.Error())
		return
	}

	defer r.Body.Close()

	if r.StatusCode != http.StatusOK {
		fmt.Errorf("Incorrect status code returned %d", r.StatusCode)
		return
	}

	var lst Listing
	contents, _ := ioutil.ReadAll(r.Body)
	decodeErr := json.Unmarshal([]byte(string(contents)), &lst)

	if decodeErr != nil {
		fmt.Errorf("Unable to decode the JSON feed: %s", decodeErr.Error())
		return
	}

	// cleanList(&lst)
	after, music = lst.Data.After, makeMap(lst)

	return
}

func makeMap(lst Listing) (music map[string]*Children) {
	music = make(map[string]*Children)

	for _, el := range lst.Data.Children {
		music[el.Data.ID] = &Children{"", el.Data}
	}

	return
}

// cleanList Removes non-youtube items because we cannot, currently, download them.
func cleanList(lst *Listing) {
	whitelistedDomains := map[string]bool{
		"youtu.be":    true,
		"youtube.com": true,
	}

	l := lst.Data.Children[0:]

	for i, el := range l {
		if !whitelistedDomains[el.Data.Domain] {
			lst.Data.Children = append(lst.Data.Children[:i], lst.Data.Children[i+1:]...)
		}
	}
}
