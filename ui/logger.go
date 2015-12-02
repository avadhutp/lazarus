package ui

import (
	"github.com/Sirupsen/logrus"
	"io"
)

var Log = &Logger{logrus.New()}

type Logger struct {
	l *logrus.Logger
}

func (l *Logger) SetOutput(w io.Writer) {
	l.SetOutput(w)
}

func (l *Logger) Error(args ...interface{}) {
	l.Error(args...)
}
