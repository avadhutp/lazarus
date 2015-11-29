package geddit

type Listing struct {
	Kind string `json:"kind"`
	Data struct {
		After    string     `json:"after"`
		Children []Children `json:"children"`
	} `json:"data"`
}

type Children struct {
	Kind string `json:"kind"`
	Data struct {
		Domain string `json:"domain"`
		Url    string `json:"url"`
		Title  string `json:"title"`
		Genre  string `json:"link_flair_text"`
		Id     string `json:"id"`
		Played bool
	} `json:"data"`
}
