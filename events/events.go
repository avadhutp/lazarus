package events

import (
	"github.com/avadhutp/lazarus/geddit"
	"github.com/gizak/termui"
)

var (
	FinishedGedditDownload = "/geddit/download/finished"
)

func FireFinishedGedditDownload(lst *geddit.Listing) {
	termui.SendCustomEvt(FinishedGedditDownload, lst)
}
