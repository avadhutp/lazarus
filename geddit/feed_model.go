package geddit

const (
	// IsDownloading The song is being downloaded
	IsDownloading = 1
	// Downloaded The song has been downloaded
	Downloaded = 2
	// NotDownloaded Cannot download the song due to some errors
	NotDownloaded = 3
	// Playing The song is being played currently
	Playing = 4
	// IsPlayed The song has finished playing
	IsPlayed = 5
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

// IsDownloading Set the status of the song to being downloaded
func (c *Children) IsDownloading() {
	c.Data.Status = IsDownloading
}

// Downloaded Set the status of the song as being downloaded
func (c *Children) Downloaded() {
	c.Data.Status = Downloaded
}

// CannotDownload Set the sthe status of the song as count not download, used when there's some error with youtube-dl
func (c *Children) CannotDownload() {
	c.Data.Status = NotDownloaded
}

// IsPlaying Set the status of the song as being currently played
func (c *Children) IsPlaying() {
	c.Data.Status = Playing
}

// FinishedPlaying Set the status of the song as finished playing
func (c *Children) FinishedPlaying() {
	c.Data.Status = IsPlayed
}
