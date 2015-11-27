package ui

import (
	"fmt"
	"github.com/avadhutp/lazarus/geddit"
	"os/exec"
	"time"
)

const (
	cmd = " youtube-dl --extract-audio -o '/tmp/lazarus/%(title)s.%(ext)s' --audio-format mp3 %s"
)

func downloadSong(el geddit.Children) {
	dlCmd := fmt.Sprintf(cmd, el.Data.Url)
	exec.Command(dlCmd)
	time.Sleep(1 * time.Second)
}
