package ui

import (
	"encoding/json"

	"github.com/stretchr/testify/assert"

	"github.com/gizak/termui"

	"github.com/avadhutp/lazarus/geddit"

	"testing"
)

var testStr = `
	{
		"kind" : "Listing",
		"data" : {
			"after" : "test-after",
			"children" : [
				{
					"kind" : "t3",
					"data" : {
						"status" : 1,
						"title" : "Test song 1",
						"id" : "1"
					}	
				},
				{
					"kind" : "t3",
					"data" : {
						"status" : 2,
						"title" : "Test song 2",
						"id" : "2"
					}	
				},
				{
					"kind" : "t3",
					"data" : {
						"status" : 3,
						"title" : "Test song 3",
						"id" : "3"
					}	
				},
				{
					"kind" : "t3",
					"data" : {
						"status" : 4,
						"title" : "Test song 4",
						"id" : "4"
					}	
				},
				{
					"kind" : "t3",
					"data" : {
						"status" : 5,
						"title" : "Test song 5",
						"id" : "5"
					}	
				}
			]
		}
	}
`

func lstToMusic(lst geddit.Listing) (music map[string]*geddit.Children) {
	music = make(map[string]*geddit.Children)

	for _, el := range lst.Data.Children {
		music[el.Data.ID] = &geddit.Children{"", el.Data}
	}

	return
}

func TestPaintSongList(t *testing.T) {
	oldRefresh := Refresh
	defer func() { Refresh = oldRefresh }()

	var lst geddit.Listing
	json.Unmarshal([]byte(testStr), &lst)

	player := Player{}
	player.Music = lstToMusic(lst)

	evt := new(termui.Event)
	evt.Data = player

	refreshCalled := false
	Refresh = func() {
		refreshCalled = true
	}

	paintSongList(*evt)

	var expectedList = []string{
		"|[»](fg-white,)|[Test song 1](fg-white,)",
		"|[■](fg-green,)|[Test song 2](fg-green,)",
		"|[✖](fg-red,)|[Test song 3](fg-red,)",
		"|[►](fg-white,bg-green)|[Test song 4](fg-white,bg-green)",
		"|[✔](fg-blue,)|[Test song 5](fg-blue,)",
	}

	assert.Equal(t, expectedList, Songs.Items)
	assert.Equal(t, len(player.Music)+songsPadding, Songs.GetHeight())
	assert.True(t, refreshCalled, "Refresh should be called after the Songs widget has been updated")
}
