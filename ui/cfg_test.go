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
		playerCmd    string
		osStatErr    error
		osIsNotExist bool
		errMsg       string
		allOkIsOk    bool
		msg          string
	}{
		{
			tmpLocation:  "",
			playerCmd:    "something",
			osStatErr:    nil,
			osIsNotExist: false,
			errMsg:       "Missing directive in the ini file",
			allOkIsOk:    false,
			msg:          "Empty tmp location indicates that the ini file did not contain tmp_location directive",
		},
		{
			tmpLocation:  "something",
			playerCmd:    "something",
			osStatErr:    errors.New("File does not exist error"),
			osIsNotExist: true,
			errMsg:       "File does not exist error",
			allOkIsOk:    false,
			msg:          "File does not exist; bubble up the encountered error",
		},
		{
			tmpLocation:  "something",
			playerCmd:    "something",
			osStatErr:    nil,
			osIsNotExist: false,
			errMsg:       "",
			allOkIsOk:    true,
			msg:          "File exists and therefore we should see no errors",
		},
		{
			tmpLocation:  "something",
			playerCmd:    "",
			osStatErr:    nil,
			osIsNotExist: false,
			errMsg:       "Missing directive in the ini file: player_cmd",
			allOkIsOk:    false,
			msg:          "Empty player_cmd should trigger an error",
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
		sut.PlayerCmd = test.playerCmd

		actual := sut.AllOk()

		switch test.allOkIsOk {
		case true:
			assert.Nil(t, actual, test.msg)
		case false:
			assert.Contains(t, actual.Error(), test.errMsg, test.msg)
		}
	}

}
