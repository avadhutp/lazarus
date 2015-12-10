package ui

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/avadhutp/lazarus/geddit"

	"github.com/stretchr/testify/assert"
)

type TestFileInfo struct {
	name string
}

func (t *TestFileInfo) Name() string {
	return t.name
}

func (t *TestFileInfo) Size() int64 {
	return 0
}

func (t *TestFileInfo) Mode() os.FileMode {
	return os.ModeDir
}

func (t *TestFileInfo) ModTime() time.Time {
	return time.Now()
}

func (t *TestFileInfo) IsDir() bool {
	return false
}

func (t *TestFileInfo) Sys() interface{} {
	return false
}

func TestNewPlayer(t *testing.T) {
	cfg := &Cfg{}
	cfg.PlayerCmd = "testcmd testargs"

	p := NewPlayer(cfg)

	assert.Equal(t, cfg, p.cfg)
	assert.Equal(t, "testcmd", p.playerCmd)
	assert.Equal(t, []string{"testargs"}, p.playerArgs)
}

func TestExpandYoutubeURL(t *testing.T) {
	tests := []struct {
		in       string
		expected string
		msg      string
	}{
		{
			in:       "https://youtu.be/zKVbJuhcze8",
			expected: "http://www.youtube.com/watch?v=zKVbJuhcze8",
			msg:      "Shortened youtu.be url should be expanded correctly",
		},
		{
			in:       "http://[fe80::%31%25en0]:8080/",
			expected: "http://[fe80::%31%25en0]:8080/",
			msg:      "Invalid url fails url.Parse and therefore should be returned as is",
		},
		{
			in:       "http://www.soundcloud.com/test",
			expected: "http://www.soundcloud.com/test",
			msg:      "Not a youtube URL, so should be returned as is",
		},
	}

	for _, test := range tests {
		actual := expandYoutubeURL(test.in)
		assert.Equal(t, test.expected, actual, test.msg)
	}
}

func TestPlayerSkip(t *testing.T) {
	oldPKill := pKill
	defer func() { pKill = oldPKill }()

	pKillCalled := false
	pKill = func(p *os.Process) error {
		fmt.Println(fmt.Sprintf("Process: %d", p.Pid))
		pKillCalled = true
		return nil
	}

	p := &os.Process{}
	testCmd := &exec.Cmd{}
	testCmd.Process = p

	tests := []struct {
		initialVal      *exec.Cmd
		shouldCallPKill bool
		msg             string
	}{
		{
			initialVal:      nil,
			shouldCallPKill: false,
			msg:             "No current song, so do not call pKill",
		},
		{
			initialVal:      testCmd,
			shouldCallPKill: true,
			msg:             "Current song playing, so should call pKill",
		},
	}

	for _, test := range tests {
		sut := Player{}
		sut.currSong = test.initialVal
		sut.Skip()

		assert.Equal(t, test.shouldCallPKill, pKillCalled, test.msg)
	}
}

func TestPlayerDownload(t *testing.T) {
	data := geddit.ChildData{
		Domain: "youtube.com",
		URL:    "http://www.youtube.com/",
		Title:  "Test song title",
		Genre:  "Hip-hop",
		ID:     "12345",
	}
	el := geddit.Children{
		Kind: "T3",
		Data: data,
	}

	oldExecCommand := execCommand
	oldCmdRun := cmdRun
	oldLogError := logError
	oldIoutilReaddir := ioutilReaddir
	oldTermuiSendCustomEvt := termuiSendCustomEvt
	defer func() {
		execCommand = oldExecCommand
		cmdRun = oldCmdRun
		logError = oldLogError
		ioutilReaddir = oldIoutilReaddir
		termuiSendCustomEvt = oldTermuiSendCustomEvt
	}()

	cmd := &exec.Cmd{}
	execCommandCalled := false
	execCommand = func(command string, args ...string) *exec.Cmd {
		execCommandCalled = true
		return cmd
	}

	cmdRunCalled := false

	logErrorCalled := false
	logError = func(msg ...interface{}) {
		logErrorCalled = true
	}

	termuiSendCustomEvtCalled := []string{}
	termuiSendCustomEvt = func(evt string, i interface{}) {
		termuiSendCustomEvtCalled = append(termuiSendCustomEvtCalled, "called")
	}

	f := &TestFileInfo{name: "12345.mp3"}
	ioutilReaddir = func(dir string) (files []os.FileInfo, err error) {
		files = append(files, f)
		return
	}

	cfg := &Cfg{}
	cfg.TmpLocation = "/tmp/location/"

	tests := []struct {
		downloadErr      error
		expectedStatus   int
		expectedFileLoc  string
		isLogErrorCalled bool
		msg              string
	}{
		{
			downloadErr:      nil,
			expectedStatus:   geddit.Downloaded,
			expectedFileLoc:  "/tmp/location/12345.mp3",
			isLogErrorCalled: false,
			msg:              "Download was successfull, so status and fileLoc should be appropriately set",
		},
		{
			downloadErr:      errors.New("Test error"),
			expectedStatus:   geddit.NotDownloaded,
			expectedFileLoc:  "",
			isLogErrorCalled: true,
			msg:              "Download was unsuccessfull, so status and fileLoc should be appropriately set",
		},
	}

	for _, test := range tests {
		execCommandCalled = false
		cmdRunCalled = false
		termuiSendCustomEvtCalled = []string{}
		logErrorCalled = false

		el.Data.FileLoc = ""
		el.Data.Status = 0

		cmdRun = func(c *exec.Cmd) error {
			cmdRunCalled = true
			return test.downloadErr
		}

		p := &Player{}
		p.cfg = cfg

		p.download(&el)

		assert.Equal(t, test.expectedStatus, el.Data.Status, test.msg)
		assert.Equal(t, test.expectedFileLoc, el.Data.FileLoc, test.msg)
		assert.True(t, execCommandCalled, "Should always be called")
		assert.True(t, cmdRunCalled, "Should always be called")
		assert.Len(t, termuiSendCustomEvtCalled, 2, "Should be called twice to update the player before and after attempting download")
		assert.Equal(t, test.isLogErrorCalled, logErrorCalled)
	}
}
