package geddit

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const (
	hots = "https://www.reddit.com/r/listentothis/hot.json?sort=hot"
)

func Get() (lst Listing) {
	r, err := http.Get(hots)

	if err != nil {
		fmt.Errorf("Unable to get the subreddit %s; error faced: %s", hots, err.Error())
		return
	}

	defer r.Body.Close()

	if r.StatusCode != http.StatusOK {
		fmt.Errorf("Incorrect status code returned %d", r.StatusCode)
		return
	}

	contents, _ := ioutil.ReadAll(r.Body)
	decodeErr := json.Unmarshal([]byte(string(contents)), &lst)

	if decodeErr != nil {
		fmt.Errorf("Unable to decode the JSON feed: %s", decodeErr.Error())
		return
	}

	return
}
