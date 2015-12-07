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
