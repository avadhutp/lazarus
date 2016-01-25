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
	gedditGet     = geddit.Get
	pKill         = (*os.Process).Kill
	sleep         = time.Sleep
	execCommand   = exec.Command
	cmdRun        = (*exec.Cmd).Run
	logError      = log.Error
	ioutilReaddir = ioutil.ReadDir

	playerRestart        func(*Player)
	playerStart          func(*Player)
	playerStartDownloads func(*Player)
	playerStartPlayback  func(*Player)
	playerDownload       func(*Player, *geddit.Children)
	playerPlay           func(*Player, *geddit.Children)
)

func init() {
	playerStart = (*Player).Start
	playerRestart = (*Player).restart
	playerStartDownloads = (*Player).startDownloads
	playerStartPlayback = (*Player).startPlayback
	playerDownload = (*Player).download
	playerPlay = (*Player).play
}

const (
	waitBeforeStartingPlayback = 5 * time.Second
	waitOnDownloadingSong      = 5 * time.Second

	redditURL = "http://www.reddit.com/%s/hot.json?sort=hot&after=%s"
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
	p.after, p.Music = gedditGet(p.getRedditURL())
	FireFinishedRedditDownload(*p)

	go playerStartDownloads(p)
	sleep(waitBeforeStartingPlayback)
	go playerStartPlayback(p)
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

	playerStart(p)
}

func (p *Player) getRedditURL() string {
	return fmt.Sprintf(redditURL, subreddit, p.after)
}

func (p *Player) startPlayback() {
	for _, k := range p.GetKeys() {
		playerPlay(p, p.Music[k])
	}
	playerRestart(p)
}

func (p *Player) play(el *geddit.Children) {
	switch el.Data.Status {
	case geddit.IsPlayed:
	case geddit.NotDownloaded:
		return
	case geddit.IsDownloading:
		sleep(waitOnDownloadingSong)
		p.play(el)
	case geddit.Downloaded:
		p.runPlayCmd(el)
	}
}

func (p *Player) runPlayCmd(el *geddit.Children) {
	el.IsPlaying()
	UpdatePlayer(*p)

	args := append(p.playerArgs, el.Data.FileLoc)
	cmd := execCommand(p.playerCmd, args...)
	p.currSong = cmd
	cmdRun(cmd)

	el.FinishedPlaying()
	UpdatePlayer(*p)
}

func (p *Player) startDownloads() {
	for _, k := range p.GetKeys() {
		playerDownload(p, p.Music[k])
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
		expandYoutubeURL(el.Data.URL),
	}
	cmd := execCommand("youtube-dl", args...)
	err := cmdRun(cmd)

	if err != nil {
		logError(fmt.Sprintf("Cannot download %s; Error encountered: %s", el.Data.URL, err.Error()))
	} else {
		p.setFileName(el)
	}

	return err
}

func (p *Player) getFileLocation(el *geddit.Children) string {
	return p.cfg.TmpLocation + el.Data.ID + ".%(ext)s"
}

func (p *Player) setFileName(el *geddit.Children) {
	files, _ := ioutilReaddir(p.cfg.TmpLocation)

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

func deleteFile(loc string) {
	if err := os.Remove(loc); err != nil {
		logError(fmt.Sprintf("Cannot delete file: %s; error encountered: %s", loc, err.Error()))
	}
}
