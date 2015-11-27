package events

import (
	"github.com/avadhutp/lazarus/geddit"
	"github.com/gizak/termui"
)

var (
	FinishedGedditDownload = "/geddit/download/finished"
	StartSongDownload      = "/lazarus/song/download"
)

func FireFinishedGedditDownload(lst *geddit.Listing) {
	termui.SendCustomEvt(FinishedGedditDownload, lst)
}

func FireStartSongDownload(lst *geddit.Listing) {
	termui.SendCustomEvt(StartSongDownload, lst)
}
