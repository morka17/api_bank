package worker

import (
	"fmt"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Logger interface {
	Print(level zerolog.Level, args ...interface{})
	Debug(args ...interface{})
	Info(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})
	Fatal(args ...interface{})
}

type logger struct {
}

func NewLogger() Logger {
	return &logger{}
}

func (l *logger) Print(level zerolog.Level, args ...interface{}) {
	log.WithLevel(level).Msg(fmt.Sprint(args...))
}
func (l *logger) Debug(args ...interface{}) {
	l.Print(zerolog.DebugLevel, args...)
}
func (l *logger) Info(args ...interface{}) {
	l.Print(zerolog.InfoLevel, args...)
}
func (l *logger) Warn(args ...interface{}) {
	l.Print(zerolog.WarnLevel, args...)
}
func (l *logger) Error(args ...interface{}) {
	l.Print(zerolog.ErrorLevel, args...)
}
func (l *logger) Fatal(args ...interface{}) {
	l.Print(zerolog.FatalLevel, args...)
}
