package ui

import (
	"os/exec"
	"sort"
	"time"

	"github.com/avadhutp/lazarus/geddit"
)

// Player Datastructure to hold songs and download/play them
type Player struct {
	Music map[string]*geddit.Children
	Keys  []string
	Cfg   *Cfg
}

// GetKeys Since all the songs are held in a map, to make the order of retrieval deterministic, we set the order ourselves using this func
func (p *Player) GetKeys() []string {
	if len(p.Keys) == 0 {
		for k := range p.Music {
			p.Keys = append(p.Keys, k)
		}

		sort.Strings(p.Keys)
	}

	return p.Keys
}

// Start Initiates the song download process
func (p *Player) Start() {
	for _, k := range p.GetKeys() {
		p.download(p.Music[k])
	}
}

func (p *Player) download(el *geddit.Children) {
	el.IsDownloading()
	UpdatePlayer(*p)

	time.Sleep(1 * time.Second)

	el.Downloaded()
	UpdatePlayer(*p)
}

func downloadSong(el geddit.Children) {
	args := []string{
		"--extract-audio",
		"-o",
		"/tmp/lazarus/" + el.Data.ID + ".mp3",
		"--audio-format",
		"mp3",
		el.Data.URL,
	}
	cmd := exec.Command("youtube-dl", args...)

	cmd.Run()
}
