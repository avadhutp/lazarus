package ui

import "github.com/gizak/termui"

const (
	finishedRedditDownload = "/lazarus/reddit/download/done"
	songListUpdated        = "/lazarus/update/songlist"
)

var (
	termuiSendCustomEvt = termui.SendCustomEvt
	termuiHandle        = termui.Handle
	termuiStopLoop      = termui.StopLoop
)

// UpdatePlayer Fires the event to redraw the player
func UpdatePlayer(player Player) {
	termuiSendCustomEvt(songListUpdated, player)
}

// FireFinishedRedditDownload Fires the event to begin the playlist display/download/play process
func FireFinishedRedditDownload(player Player) {
	termuiSendCustomEvt(finishedRedditDownload, player)
}

// EventHandler Registers all the event handlers
func EventHandler() {
	termuiHandle("/sys/kbd/q", func(termui.Event) { termuiStopLoop() })

	termuiHandle(finishedRedditDownload, paintSongList)
	termuiHandle(songListUpdated, paintSongList)
}

// PlayerControlEventHandler Adds keyboard controls to control the player's playback
func PlayerControlEventHandler(player PlayerInterface) {
	termuiHandle("/sys/kbd/s", func(termui.Event) { player.Skip() })
}
