package ui

import (
	"fmt"
	"net/url"
	"os/exec"
	"sort"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"

	"github.com/avadhutp/lazarus/geddit"
)

const (
	waitBeforeStartingPlayback = 5 * time.Second
	waitOnDownloadingSong      = 5 * time.Second
)

// NewPlayer Constructs a new player object with the pre-requisites
func NewPlayer(cfg *Cfg) Player {
	p := Player{map[string]*geddit.Children{}, []string{}, cfg}

	return p
}

// Player Datastructure to hold songs and download/play them
type Player struct {
	Music map[string]*geddit.Children
	Keys  []string
	Cfg   *Cfg
}

// Start Re/starts the entire download & play cycle when called; will generally be issued in main() or when the current *Player.startPlayback() loop is done
func (p *Player) Start(rURL string) {
	lst := geddit.Get(rURL)
	p.Music = lst
	FireFinishedRedditDownload(*p)

	go p.startDownloads()
	time.Sleep(waitBeforeStartingPlayback)
	go p.startPlayback()
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

func (p *Player) startPlayback() {
	for _, k := range p.GetKeys() {
		p.play(p.Music[k])
	}
}

func (p *Player) play(el *geddit.Children) {
	switch el.Data.Status {
	case geddit.IsPlayed:
	case geddit.NotDownloaded:
		return
	case geddit.IsDownloading:
		time.Sleep(waitOnDownloadingSong)
		p.play(el)
	case geddit.Downloaded:
		p.runPlayCmd(el)
	}
}

func (p *Player) runPlayCmd(el *geddit.Children) {
	el.IsPlaying()
	UpdatePlayer(*p)

	args := []string{
		"--play-and-exit",
		p.getFileLocation(el),
	}

	cmd := exec.Command("cvlc", args...)
	cmd.Run()

	el.FinishedPlaying()
	UpdatePlayer(*p)
}

func (p *Player) startDownloads() {
	for _, k := range p.GetKeys() {
		p.download(p.Music[k])
	}
}

func (p *Player) download(el *geddit.Children) {
	el.IsDownloading()
	UpdatePlayer(*p)

	switch p.runDownloadCmd(el) {
	case nil:
		el.Downloaded()
	default:
		el.CannotDownload()
	}

	UpdatePlayer(*p)
}

func (p *Player) runDownloadCmd(el *geddit.Children) error {
	args := []string{
		"--extract-audio",
		"-o",
		p.getFileLocation(el),
		"--audio-format",
		"mp3",
		expandYoutubeURL(el.Data.URL),
	}
	cmd := exec.Command("youtube-dl", args...)
	err := cmd.Run()

	if err != nil {
		log.Error(fmt.Sprintf("Cannot download %s; Error encountered: %s", el.Data.URL, err.Error()))
	}

	return err
}

func (p *Player) getFileLocation(el *geddit.Children) string {
	return p.Cfg.TmpLocation + el.Data.ID + ".mp3"
}

func expandYoutubeURL(URL string) string {
	u, err := url.Parse(URL)

	if err != nil {
		return URL
	}

	if u.Host == "youtu.be" {
		return fmt.Sprintf("http://www.youtube.com/watch?v=%s", strings.TrimPrefix(u.Path, "/"))
	}

	return URL
}
