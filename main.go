package main

import (
	"github.com/gizak/termui"
	// "fmt"
	"github.com/avadhutp/lazarus/geddit"
	"github.com/avadhutp/lazarus/ui"
)

func main() {
	lst := geddit.Get()
	songsWidget := ui.SongsWidget(lst)

	err := termui.Init()
	if err != nil {
		panic(err)
	}
	defer termui.Close()

	termui.Render(songsWidget)
	termui.Handle("/sys/kbd/q", func(termui.Event) {
		termui.StopLoop()
	})
	termui.Loop()
}
