package stdlog

import (
	"fmt"
	"time"
)

const CURRENT_LEVEL = DEBUG

const (
	DEBUG = iota
	INFO
	ERROR
	FATAL
)

var level_label = []string{"DBG", "INF", "ERR", "FAT"}

func log(level int, msg string) {

	if level < DEBUG || level > FATAL {
		level = ERROR
		logf(FATAL, "invalid loglevel '%v' passed to stdlog. Logging at UNKNOWN level", level)
	}

	now := time.Now().Local()
	ts := now.Format("2006-01-02|15:04:05.0000000")

	if level >= CURRENT_LEVEL {
		fmt.Printf("%v|%v|%v\n", ts, level_label[level], msg)
	}

}

func logf(level int, format string, a ...interface{}) {
	log(level, fmt.Sprintf(format, a...))
}

func Debug(msg string) {
	log(DEBUG, msg)
}

func Info(msg string) {
	log(INFO, msg)
}

func Error(msg string) {
	log(ERROR, msg)
}

func Fatal(msg string) {
	log(FATAL, msg)
}

func Debugf(format string, a ...interface{}) {
	log(DEBUG, fmt.Sprintf(format, a...))
}

func Infof(format string, a ...interface{}) {
	log(INFO, fmt.Sprintf(format, a...))
}

func Errorf(format string, a ...interface{}) {
	log(ERROR, fmt.Sprintf(format, a...))
}

func Fatalf(format string, a ...interface{}) {
	log(FATAL, fmt.Sprintf(format, a...))
}
