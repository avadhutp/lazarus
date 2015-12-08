package ui

import (
	"testing"

	"github.com/gizak/termui"

	"github.com/stretchr/testify/assert"
)

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
	defer func() { termuiHandle = oldTermuiHandle }()

	callList := []string{}
	termuiHandle = func(evt string, f func(termui.Event)) {
		callList = append(callList, evt)
	}

	EventHandler()

	assert.Contains(t, callList, "/sys/kbd/q", "Event handler for keypress q added")
	assert.Contains(t, callList, finishedRedditDownload, "Event handler for fiinshedRedditDownload added")
	assert.Contains(t, callList, songListUpdated, "Event handler for songListUpdated added")
}
