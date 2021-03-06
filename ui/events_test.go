package ui

import (
	"testing"

	"github.com/gizak/termui"
	"github.com/stretchr/testify/assert"
)

type MockPlayer struct {
	skipCalled bool
}

func (m *MockPlayer) Skip()             { m.skipCalled = true }
func (m *MockPlayer) Start()            {}
func (m *MockPlayer) GetKeys() []string { return nil }

func TestEventTriggers(t *testing.T) {
	oldTermuiSendCustomEvt := termuiSendCustomEvt
	defer func() { termuiSendCustomEvt = oldTermuiSendCustomEvt }()

	var actualData Player
	var actualEvt string
	termuiSendCustomEvt = func(evt string, data interface{}) {
		actualData = data.(Player)
		actualEvt = evt
	}

	tests := []struct {
		trigger     func(Player)
		expectedEvt string
		msg         string
	}{
		{
			trigger:     UpdatePlayer,
			expectedEvt: songListUpdated,
			msg:         "Should be triggered with the songListUpdated event path",
		},
		{
			trigger:     FireFinishedRedditDownload,
			expectedEvt: finishedRedditDownload,
			msg:         "Should be triggered with the finishedRedditDownload event path",
		},
	}

	p := Player{}

	for _, test := range tests {
		test.trigger(p)
		assert.Equal(t, p, actualData, "Should pass the player obj")
		assert.Equal(t, test.expectedEvt, actualEvt, test.msg)
	}
}

func TestEventHandler(t *testing.T) {
	oldTermuiHandle := termuiHandle
	defer func() {
		termuiHandle = oldTermuiHandle
	}()

	callList := map[string]func(termui.Event){}
	termuiHandle = func(evt string, f func(termui.Event)) {
		callList[evt] = f
	}

	EventHandler()

	assert.Contains(t, callList, finishedRedditDownload, "Event handler for fiinshedRedditDownload added")
	assert.Contains(t, callList, songListUpdated, "Event handler for songListUpdated added")
}

func TestPlayerControlEventHandler(t *testing.T) {
	oldTermuiHandle := termuiHandle
	oldTermuiStopLoop := termuiStopLoop
	defer func() {
		termuiHandle = oldTermuiHandle
		termuiStopLoop = oldTermuiStopLoop
	}()

	callList := map[string]func(termui.Event){}
	termuiHandle = func(evt string, f func(termui.Event)) {
		callList[evt] = f
	}

	stopLoopCalled := false
	termuiStopLoop = func() {
		stopLoopCalled = true
	}

	sut := new(MockPlayer)

	PlayerControlEventHandler(sut)

	assert.Contains(t, callList, "/sys/kbd/s", "Event handler for keypress s (skip) added")
	assert.Contains(t, callList, "/sys/kbd/q", "Event handler for keypress q (quit) added")

	callList["/sys/kbd/s"](termui.Event{})
	assert.True(t, sut.skipCalled, "Skip should be called when the event is fired")

	callList["/sys/kbd/q"](termui.Event{})
	assert.True(t, stopLoopCalled, "Stop loop is called when the *q* key is pressed")
}
