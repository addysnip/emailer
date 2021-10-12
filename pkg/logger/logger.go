package logger

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/getsentry/sentry-go"
)

const (
	LOG4G_VERSION = "log4g-v0.99.0"
	LOG4G_MAJOR   = 0
	LOG4G_MINOR   = 99
	LOG4G_PATCH   = 0
)

type Level int

const (
	DEBUG Level = iota
	INFO
	WARNING
	ERROR
	CRITICAL
	FATAL
)

var (
	levelStrings = [...]string{"DEBUG", "INFO", "WARN", "ERROR", "CRITICAL", "FATAL"}
)

var logLevel = INFO
var OutputStream = os.Stdout
var ErrorStream = os.Stderr

var sentryReady = false

type Logger struct {
	Date     time.Time
	Category string
	Level    Level
	Message  string
}

func InitSentry(dsn string) error {
	err := sentry.Init(sentry.ClientOptions{Dsn: dsn})
	if err != nil {
		return err
	}
	sentryReady = true
	return nil
}

func SetLogLevel(level Level) (bool, error) {
	if int(level) > len(levelStrings) || level < 0 {
		return false, errors.New("invalid log level")
	}

	logLevel = level
	return true, nil
}

func Category(category string) *Logger {
	logger := Logger{Category: category}

	return &logger
}

func (l *Logger) write() {
	if int(logLevel) <= int(l.Level) {
		msg := fmt.Sprintf("[%s] [%s] %-8s - %s", l.Date.Format(time.RFC3339), l.Category, l.Level, l.Message)

		if int(ERROR) <= int(l.Level) {
			fmt.Fprintln(ErrorStream, msg)
		} else {
			fmt.Fprintln(OutputStream, msg)
		}
	}
}

func (l *Logger) handle(level Level, format string, args ...interface{}) {
	l.Date = time.Now()
	l.Level = level
	l.Message = fmt.Sprintf(format, args...)
	l.write()
}

func (l *Logger) Debug(format string, args ...interface{}) {
	l.handle(DEBUG, format, args...)
}

func (l *Logger) Info(format string, args ...interface{}) {
	l.handle(INFO, format, args...)
}

func (l *Logger) Warning(format string, args ...interface{}) {
	l.handle(WARNING, format, args...)
}

func (l *Logger) Error(format string, args ...interface{}) {
	if sentryReady {
		sentry.CaptureException(fmt.Errorf(format, args...))
	}
	l.handle(ERROR, format, args...)
}

func (l *Logger) Critical(format string, args ...interface{}) {
	if sentryReady {
		sentry.CaptureException(fmt.Errorf(format, args...))
	}
	l.handle(CRITICAL, format, args...)
}

func (l *Logger) Fatal(format string, args ...interface{}) {
	if sentryReady {
		sentry.CaptureException(fmt.Errorf(format, args...))
	}
	l.handle(FATAL, format, args...)
	os.Exit(1)
}

func (l Level) String() string {
	return levelStrings[l]
}

func (l Level) Index() int {
	return int(l)
}
