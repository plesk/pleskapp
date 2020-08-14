// Copyright 1999-2020. Plesk International GmbH.

package utils

import (
	"log"
	"os"

	"github.com/plesk/pleskapp/plesk/locales"
)

var Log *logger

type logger struct {
	stdout *log.Logger
	stderr *log.Logger
	level  int
}

func (l *logger) SetLevel(lv int) {
	Log.level = lv
}

func (l *logger) HasDebug() bool {
	return l.level >= 2
}

func (l *logger) PrintL(s string, vars ...interface{}) {
	l.stdout.Print(locales.L.Get(s, vars...))
}

func (l *logger) Print(s string) {
	l.stdout.Print(s)
}

func (l *logger) Verbose(s string) {
	if l.level >= 1 {
		l.stderr.Print(s)
	}
}

func (l *logger) Debug(s string) {
	if l.level >= 2 {
		l.stderr.Print(s)
	}
}

func (l *logger) Error(s string) {
	l.stderr.Print(s)
}

func init() {
	Log = &logger{
		stdout: log.New(os.Stdout, "", 0),
		stderr: log.New(os.Stderr, "", 0),
		level:  0,
	}
}
