package ui

import (
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func getSUT() *Cfg {
	return &Cfg{}
}

func TestAllOk(t *testing.T) {
	oldOsStat := osStat
	oldOsIstNotExist := osIsNotExist
	defer func() {
		osStat = oldOsStat
		osIsNotExist = oldOsIstNotExist
	}()

	tests := []struct {
		tmpLocation  string
		osStatErr    error
		osIsNotExist bool
		errMsg       string
		msg          string
	}{
		{
			tmpLocation:  "",
			osStatErr:    nil,
			osIsNotExist: false,
			errMsg:       "Missing directive in the ini file",
			msg:          "Empty tmp location indicates that the ini file did not contain tmp_location directive",
		},
		{
			tmpLocation:  "something",
			osStatErr:    errors.New("File does not exist error"),
			osIsNotExist: true,
			errMsg:       "File does not exist error",
			msg:          "File does not exist; bubble up the encountered error",
		},
	}

	for _, test := range tests {
		osStat = func(name string) (os.FileInfo, error) {
			return nil, test.osStatErr
		}
		osIsNotExist = func(err error) bool {
			return test.osIsNotExist
		}
		sut := getSUT()

		sut.TmpLocation = test.tmpLocation

		actual := sut.AllOk()

		assert.Contains(t, actual.Error(), test.errMsg, test.msg)
	}

}
