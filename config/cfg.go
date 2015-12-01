package common

// Config Maps to the lazarus config ini file.
type Cfg struct {
	TmpLocation string `ini:"tmp_location"`
	MazSize     string `ini:"max_size"`
	LogFile     string `ini:"log_file"`
}

func (c *Cfg) Kanah() (bool, error) {
	return true, nil
}
