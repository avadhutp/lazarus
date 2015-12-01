package ui

import (
	"os/exec"
	"sort"
	"time"

	"github.com/avadhutp/lazarus/geddit"
)

type Player struct {
	Music    map[string]*geddit.Children
	Keys     []string
	Cfg_main *Cfg
}

func (p *Player) GetKeys() []string {
	if len(p.Keys) == 0 {
		for k, _ := range p.Music {
			p.Keys = append(p.Keys, k)
		}

		sort.Strings(p.Keys)
	}

	return p.Keys
}

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
		"/tmp/lazarus/" + el.Data.Id + ".mp3",
		"--audio-format",
		"mp3",
		el.Data.Url,
	}
	cmd := exec.Command("youtube-dl", args...)

	cmd.Run()
}
