package ui

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func getSUT() *Cfg {
	return &Cfg{}
}

func TestAllOk(t *testing.T) {
	sut := getSUT()

	sut.TmpLocation = ""

	actual := sut.AllOk()

	assert.Contains(t, actual.Error(), "Missing directive in the ini file")
}
