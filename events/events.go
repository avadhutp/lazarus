package events

import (
	"github.com/avadhutp/lazarus/geddit"
	"github.com/gizak/termui"
)

var (
	FinishedGedditDownload = "/geddit/download/finished"
	DoneSongDownload       = "/lazarus/song/download"
)

func FireFinishedGedditDownload(music map[string]geddit.Children) {
	termui.SendCustomEvt(FinishedGedditDownload, music)
}

func EventHandler() {
	termui.Handle("/sys/kbd/q", func(termui.Event) { termui.StopLoop() })
	termui.Handle(events.FinishedGedditDownload, updateSongList)
}

type Player struct {
	Music map[string]*geddit.Children
}

func (p *Player) Start() {
	// for _, el := range p.lst {

	// }
}
