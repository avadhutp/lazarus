package ui

import (
	"os"

	"errors"
)

// Config Maps to the lazarus config ini file.
type Cfg struct {
	TmpLocation string `ini:"tmp_location"`
	MazSize     string `ini:"max_size"`
	LogFile     string `ini:"log_file"`
}

// AllOk Makes sure that the mandatory values in the config file exist and are ok.
func (c *Cfg) AllOk() error {
	if err := c.isTmpLocationOk(); err != nil {
		return err
	}

	return nil
}

// isTmpLocationOk Ensures that the tmp_location from the ini file exists.
func (c *Cfg) isTmpLocationOk() error {
	if c.TmpLocation == "" {
		return errors.New("Missing directive in ini file: tmp_location")
	}

	if _, err := os.Stat(c.TmpLocation); os.IsNotExist(err) {
		return err
	}

	return nil
}
