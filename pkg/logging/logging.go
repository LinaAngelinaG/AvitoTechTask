package logging

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"path"
	"runtime"
)

var e *logrus.Entry

type Logger struct {
	*logrus.Entry
}

func GetLogger() Logger {
	return Logger{e}
}

func (logger *Logger) GetLoggerWithField(k string, v interface{}) Logger {
	return Logger{logger.WithField(k, v)}
}

type hookWriter struct {
	Writer    []io.Writer
	LogLevels []logrus.Level
}

func (hook *hookWriter) Fire(entry *logrus.Entry) error {
	line, err := entry.String()
	if err != nil {
		return err
	}
	for _, w := range hook.Writer {
		w.Write([]byte(line))
	}
	return err
}

func (hook *hookWriter) Levels() []logrus.Level {
	return hook.LogLevels
}

func init() {
	logger := logrus.New()
	logger.SetReportCaller(true)
	logger.Formatter = &logrus.TextFormatter{
		CallerPrettyfier: func(frame *runtime.Frame) (function string, file string) {
			filename := path.Base(frame.File)
			return fmt.Sprintf("%s()", frame.Function), fmt.Sprintf("%s:%d", filename, frame.Line)
		},
		FullTimestamp: true,
		DisableColors: false,
	}
	err := os.MkdirAll("logs", 0777)
	if err != nil {
		panic(err)
	}
	files, err := os.OpenFile("logs/all.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0640)
	if err != nil {
		panic(err)
	}

	logger.SetOutput(io.Discard)
	logger.AddHook(&hookWriter{
		Writer:    []io.Writer{files, os.Stdout},
		LogLevels: logrus.AllLevels,
	})

	logger.SetLevel(logrus.TraceLevel)
	e = logrus.NewEntry(logger)
}
