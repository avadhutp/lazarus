package ui

import (
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"

	"github.com/avadhutp/lazarus/geddit"
)

// DI
var (
	pKill = (*os.Process).Kill
)

const (
	waitBeforeStartingPlayback = 5 * time.Second
	waitOnDownloadingSong      = 5 * time.Second

	redditURL = "https://www.reddit.com/%s/hot.json?sort=hot&after=%s"
	subreddit = "r/listentothis"
)

// PlayerInterface Makes testing easier
type PlayerInterface interface {
	Start()
	Skip()
	GetKeys() []string
}

// NewPlayer Constructs a new player object with the pre-requisites
func NewPlayer(cfg *Cfg) Player {
	p := Player{}
	p.cfg = cfg

	playerCmd := strings.Split(cfg.PlayerCmd, " ")
	p.playerCmd = playerCmd[:1][0]
	p.playerArgs = playerCmd[1:]

	return p
}

// Player Datastructure to hold songs and download/play them
type Player struct {
	Music      map[string]*geddit.Children
	after      string
	keys       []string
	cfg        *Cfg
	currSong   *exec.Cmd
	playerCmd  string
	playerArgs []string
}

// Start Re/starts the entire download & play cycle when called; will generally be issued in main() or when the current *Player.startPlayback() loop is done
func (p *Player) Start() {
	p.after, p.Music = geddit.Get(p.getRedditURL())
	FireFinishedRedditDownload(*p)

	go p.startDownloads()
	time.Sleep(waitBeforeStartingPlayback)
	go p.startPlayback()
}

// Skip skips the currently playing song
func (p *Player) Skip() {
	if p.currSong != nil {
		pKill(p.currSong.Process)
	}
}

// GetKeys Since all the songs are held in a map, to make the order of retrieval deterministic, we set the order ourselves using this func
func (p *Player) GetKeys() []string {
	if len(p.keys) == 0 {
		for k := range p.Music {
			p.keys = append(p.keys, k)
		}

		sort.Strings(p.keys)
	}

	return p.keys
}

func (p *Player) restart() {
	p.Music = make(map[string]*geddit.Children)
	p.keys = make([]string, 0)

	p.Start()
}

func (p *Player) getRedditURL() string {
	return fmt.Sprintf(redditURL, subreddit, p.after)
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

	args := append(p.playerArgs, el.Data.FileLoc)
	cmd := exec.Command(p.playerCmd, args...)
	p.currSong = cmd
	cmd.Run()

	el.FinishedPlaying()
	UpdatePlayer(*p)
}

func (p *Player) startDownloads() {
	for _, k := range p.GetKeys() {
		p.download(p.Music[k])
	}

	p.restart()
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
		expandYoutubeURL(el.Data.URL),
	}
	cmd := exec.Command("youtube-dl", args...)
	err := cmd.Run()

	if err != nil {
		log.Error(fmt.Sprintf("Cannot download %s; Error encountered: %s", el.Data.URL, err.Error()))
	} else {
		p.setFileName(el)
	}

	return err
}

func (p *Player) getFileLocation(el *geddit.Children) string {
	return p.cfg.TmpLocation + el.Data.ID + ".%(ext)s"
}

func (p *Player) setFileName(el *geddit.Children) {
	files, _ := ioutil.ReadDir(p.cfg.TmpLocation)

	for _, f := range files {
		if strings.Split(f.Name(), ".")[0] == el.Data.ID {
			el.Data.FileLoc = p.cfg.TmpLocation + f.Name()
		}
	}
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
