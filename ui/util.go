package ui

import (
	"github.com/avadhutp/lazarus/geddit"
	"os/exec"
)

func DownloadSong(el geddit.Children) {
	args := []string{
		"--extract-audio",
		"-o",
		"/tmp/lazarus/" + el.Data.Id + ".mp3",
		"--audio-format",
		"mp3",
		el.Data.Url,
	}
	cmd := exec.Command("youtube-dl", args...)

	cmd.Run()
}
