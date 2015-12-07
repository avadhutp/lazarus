package geddit

import (
	"github.com/stretchr/testify/assert"

	"testing"
)

func getSut() *Children {
	return &Children{}
}

func TestIsDownloading(t *testing.T) {
	sut := getSut()
	assert.Zero(t, sut.Data.Status)

	sut.IsDownloading()
	assert.Equal(t, IsDownloading, sut.Data.Status)
}

func TestDownloaded(t *testing.T) {
	sut := getSut()
	assert.Zero(t, sut.Data.Status)

	sut.Downloaded()
	assert.Equal(t, Downloaded, sut.Data.Status)
}

func TestNotDownloaded(t *testing.T) {
	sut := getSut()
	assert.Zero(t, sut.Data.Status)

	sut.CannotDownload()
	assert.Equal(t, NotDownloaded, sut.Data.Status)
}

func TestPlaying(t *testing.T) {
	sut := getSut()
	assert.Zero(t, sut.Data.Status)

	sut.IsPlaying()
	assert.Equal(t, Playing, sut.Data.Status)
}

func TestIsPlayed(t *testing.T) {
	sut := getSut()
	assert.Zero(t, sut.Data.Status)

	sut.FinishedPlaying()
	assert.Equal(t, IsPlayed, sut.Data.Status)
}
