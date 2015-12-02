package geddit

const (
	// IsDownloading The song is being downloaded
	IsDownloading = 1
	// Downloaded The song has been downloaded
	Downloaded = 2
	// Playing The song is being played currently
	Playing = 3
	// IsPlayed The song has finished playing
	IsPlayed = 4
)

// Listing structure mapping the json returned by reddit's API
type Listing struct {
	Kind string `json:"kind"`
	Data struct {
		After    string     `json:"after"`
		Children []Children `json:"children"`
	} `json:"data"`
}

// Children sub-structure mapping the json returned by reddit's API
type Children struct {
	Kind string `json:"kind"`
	Data struct {
		Domain string `json:"domain"`
		URL    string `json:"url"`
		Title  string `json:"title"`
		Genre  string `json:"link_flair_text"`
		ID     string `json:"id"`
		Played bool
		Status int
	} `json:"data"`
}

// IsDownloading Set the status of the song (Children) to being downloaded
func (c *Children) IsDownloading() {
	c.Data.Status = IsDownloading
}

// Downloaded Set the status of the song (Children) as being downloaded
func (c *Children) Downloaded() {
	c.Data.Status = Downloaded
}
