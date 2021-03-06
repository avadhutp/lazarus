package ui

import (
	"fmt"
	"os"
)

var (
	osStat       = os.Stat
	osIsNotExist = os.IsNotExist
)

// Cfg Maps to the lazarus config ini file.
type Cfg struct {
	TmpLocation string `ini:"tmp_location"`
	PlayerCmd   string `ini:"player_cmd"`
	MazSize     string `ini:"max_size"`
}

// AllOk Makes sure that the mandatory values in the config file exist and are ok.
func (c *Cfg) AllOk() error {
	if err := c.isTmpLocationOk(); err != nil {
		return err
	}

	if err := c.isPlayerCmdOk(); err != nil {
		return err
	}

	return nil
}

func (c *Cfg) isTmpLocationOk() error {
	return isLocationOk(c.TmpLocation, "tmp_location")
}

func (c *Cfg) isPlayerCmdOk() error {
	if c.PlayerCmd == "" {
		return fmt.Errorf("Missing directive in the ini file: player_cmd")
	}

	return nil
}

func isLocationOk(loc string, name string) error {
	if loc == "" {
		return fmt.Errorf("Missing directive in the ini file: %s", name)
	}

	if _, err := osStat(loc); osIsNotExist(err) {
		return err
	}

	return nil
}
