package ui

import (
	log "github.com/Sirupsen/logrus"
	"os/exec"
	"sort"

	"github.com/avadhutp/lazarus/geddit"
)

// NewPlayer Constructs a new player object with the pre-requisites
func NewPlayer(m map[string]*geddit.Children, cfg *Cfg) Player {
	return Player{m, []string{}, cfg}
}

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

	switch err := p.runCmd(el); err {
	case nil:
		el.Downloaded()
	default:
		log.Error("Problems Padron: %s", err.Error())
		el.CannotDownload()
	}

	UpdatePlayer(*p)
}

func (p *Player) runCmd(el *geddit.Children) error {
	args := []string{
		"--extract-audio",
		"-o",
		"/tmp/lazarus/" + el.Data.ID + ".mp3",
		"--audio-format",
		"mp3",
		el.Data.URL,
	}
	cmd := exec.Command("youtube-dl", args...)
	err := cmd.Run()

	return err
}
