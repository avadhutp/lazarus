package player

import (
	"github.com/avadhutp/lazarus/geddit"
	"github.com/gizak/termui"
	"sort"
	"time"
)

const (
	FinishedRedditDownload = "/lazarus/reddit/download/done"
	SongListUpdated        = "/lazarus/update/songlist"
)

func UpdateSongUI(playerObj Player) {
	termui.SendCustomEvt(SongListUpdated, playerObj)
}

func FireFinishedRedditDownload(playerObj Player) {
	termui.SendCustomEvt(FinishedRedditDownload, playerObj)
}

func EventHandler() {
	termui.Handle("/sys/kbd/q", func(termui.Event) { termui.StopLoop() })

	termui.Handle(FinishedRedditDownload, paintSongList)
	termui.Handle(SongListUpdated, paintSongList)
}

type Player struct {
	Music map[string]*geddit.Children
	Keys  []string
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
	UpdateSongUI(*p)

	time.Sleep(1 * time.Second)

	el.Downloaded()
	UpdateSongUI(*p)
}
