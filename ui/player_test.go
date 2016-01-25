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

func getTestMusic() geddit.Children {
	data := geddit.ChildData{
		Domain: "youtube.com",
		URL:    "http://www.youtube.com/",
		Title:  "Test song title",
		Genre:  "Hip-hop",
		ID:     "12345",
	}

	return geddit.Children{
		Kind: "T3",
		Data: data,
	}
}

func TestPlayerStart(t *testing.T) {
	oldGedditGet := gedditGet
	oldPlayerStartDownloads := playerStartDownloads
	oldPlayerStartPlayback := playerStartPlayback
	oldSleep := sleep
	oldTermuiSendCustomEvt := termuiSendCustomEvt

	defer func() {
		gedditGet = oldGedditGet
		playerStartDownloads = oldPlayerStartDownloads
		playerStartPlayback = oldPlayerStartPlayback
		sleep = oldSleep
		termuiSendCustomEvt = oldTermuiSendCustomEvt
	}()

	el := getTestMusic()
	ret := map[string]*geddit.Children{
		"12345": &el,
	}
	gedditGet = func(rURL string) (string, map[string]*geddit.Children) {
		return "after-test", ret
	}

	testChannel := make(chan bool)

	startDownloadsCalled := false
	playerStartDownloads = func(p *Player) {
		testChannel <- true
		startDownloadsCalled = true
	}

	startPlaybackCalled := false
	playerStartPlayback = func(p *Player) {
		testChannel <- true
		startPlaybackCalled = true
	}

	sleepCalled := false
	sleep = func(t time.Duration) {
		sleepCalled = true
	}

	var evtFired string
	termuiSendCustomEvt = func(evt string, data interface{}) {
		evtFired = evt
	}

	p := new(Player)
	p.Start()

	assert.Equal(t, "after-test", p.after)
	assert.Equal(t, ret, p.Music)
	assert.Equal(t, finishedRedditDownload, evtFired)
	assert.True(t, sleepCalled)

	<-testChannel
	<-testChannel
	assert.True(t, startDownloadsCalled)
	assert.True(t, startPlaybackCalled)
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
	el := getTestMusic()

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

func TestPlayerPlay(t *testing.T) {
	el := getTestMusic()

	oldExecCommand := execCommand
	oldCmdRun := cmdRun
	oldSleep := sleep
	oldTermuiSendCustomEvt := termuiSendCustomEvt
	oldDeleteFile := deleteFile
	defer func() {
		execCommand = oldExecCommand
		cmdRun = oldCmdRun
		sleep = oldSleep
		termuiSendCustomEvt = oldTermuiSendCustomEvt
		deleteFile = oldDeleteFile
	}()

	cmd := &exec.Cmd{}
	execCommandCalled := false
	execCommand = func(command string, args ...string) *exec.Cmd {
		execCommandCalled = true
		return cmd
	}

	termuiSendCustomEvtCalled := []string{}
	termuiSendCustomEvt = func(evt string, i interface{}) {
		termuiSendCustomEvtCalled = append(termuiSendCustomEvtCalled, "called")
	}

	cmdRunCalled := false
	cmdRun = func(c *exec.Cmd) error {
		cmdRunCalled = true
		return nil
	}

	deleteFileCalled := false
	deleteFile = func(loc string) {
		deleteFileCalled = true
	}

	sleepCalled := false
	sleep = func(t time.Duration) {
		el.Data.Status = geddit.Downloaded
		sleepCalled = true
	}

	tests := []struct {
		initialStatus   int
		isPlayCmdCalled bool
		isSleepCalled   bool
		msg             string
	}{
		{
			initialStatus:   geddit.IsPlayed,
			isPlayCmdCalled: false,
			isSleepCalled:   false,
			msg:             "Song is already played",
		},
		{
			initialStatus:   geddit.NotDownloaded,
			isPlayCmdCalled: false,
			isSleepCalled:   false,
			msg:             "Song download failed",
		},
		{
			initialStatus:   geddit.Downloaded,
			isPlayCmdCalled: true,
			isSleepCalled:   false,
			msg:             "Song has been downloaded, it should be played",
		},
		{
			initialStatus:   geddit.IsDownloading,
			isPlayCmdCalled: true,
			isSleepCalled:   true,
			msg:             "Song is downloading, sleep should be called and it should play on retry",
		},
	}

	for _, test := range tests {
		execCommandCalled = false
		cmdRunCalled = false
		deleteFileCalled = false
		termuiSendCustomEvtCalled = []string{}
		sleepCalled = false

		el.Data.Status = test.initialStatus
		p := &Player{}
		p.play(&el)

		switch test.isPlayCmdCalled {
		case true:
			assert.Equal(t, geddit.IsPlayed, el.Data.Status, test.msg)
			assert.Len(t, termuiSendCustomEvtCalled, 2)
			assert.True(t, execCommandCalled, test.msg)
			assert.True(t, cmdRunCalled, test.msg)
			assert.True(t, deleteFileCalled, test.msg)
		case false:
			assert.Len(t, termuiSendCustomEvtCalled, 0)
			assert.False(t, execCommandCalled, test.msg)
			assert.False(t, cmdRunCalled, test.msg)
		}
		assert.Equal(t, test.isSleepCalled, sleepCalled, test.msg)
	}
}

func TestPlayerGetRedditURL(t *testing.T) {
	p := &Player{}
	p.after = "afterhash"

	actual := p.getRedditURL()

	assert.Equal(t, "http://www.reddit.com/r/listentothis/hot.json?sort=hot&after=afterhash", actual)
}

func TestPlayerRestart(t *testing.T) {
	oldPlayerStart := playerStart
	defer func() { playerStart = oldPlayerStart }()

	el := getTestMusic()

	playerStartCalled := false
	playerStart = func(p *Player) {
		playerStartCalled = true
	}

	player := new(Player)
	player.Music = map[string]*geddit.Children{"12345": &el}
	player.keys = []string{"12345"}

	player.restart()

	assert.Empty(t, player.Music)
	assert.Empty(t, player.keys)
	assert.True(t, playerStartCalled)
}

func TestPlayerStartDownloads(t *testing.T) {
	oldPlayerDownload := playerDownload
	defer func() { playerDownload = oldPlayerDownload }()

	callList := []string{}
	playerDownload = func(p *Player, el *geddit.Children) {
		callList = append(callList, el.Data.URL)
	}

	el := getTestMusic()
	player := new(Player)
	player.Music = map[string]*geddit.Children{"12345": &el}

	player.startDownloads()

	assert.Equal(t, []string{el.Data.URL}, callList, "All songs in the music list should be sent to Player.download")
}

func TestPlayerStartPlayback(t *testing.T) {
	oldPlayerPlay := playerPlay
	oldPlayerRestart := playerRestart
	defer func() {
		playerRestart = oldPlayerRestart
		playerPlay = oldPlayerPlay
	}()

	callList := []string{}
	playerPlay = func(p *Player, el *geddit.Children) {
		callList = append(callList, el.Data.FileLoc)
	}

	playerRestartCalled := false
	playerRestart = func(p *Player) {
		playerRestartCalled = true
	}

	el := getTestMusic()
	player := new(Player)
	player.Music = map[string]*geddit.Children{"12345": &el}

	player.startPlayback()

	assert.Equal(t, []string{el.Data.FileLoc}, callList, "All songs in the music list should be sent to Player.play")
	assert.True(t, playerRestartCalled)
}

func TestPlayerDeleteFile(t *testing.T) {
	oldOSRemove := osRemove
	oldLogError := logError

	defer func() {
		osRemove = oldOSRemove
		logError = oldLogError
	}()

	logErrorCalled := false
	logError = func(...interface{}) {
		logErrorCalled = true
	}

	tests := []struct {
		err            error
		isLogErrCalled bool
		msg            string
	}{
		{
			err:            errors.New("Sample error"),
			isLogErrCalled: true,
			msg:            "Error encountered while trying to delete file",
		},
	}

	for _, test := range tests {
		osRemove = func(loc string) error {
			return test.err
		}

		deleteFile("/file/location")

		assert.Equal(t, test.isLogErrCalled, logErrorCalled, test.msg)
	}

}
