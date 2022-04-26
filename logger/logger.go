package logger

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"runtime"

	"github.com/sirupsen/logrus"
)

type writerHook struct {
	Writer    []io.Writer
	LogLevels []logrus.Level
}

func (hook *writerHook) Fire(entry *logrus.Entry) error {
	line, err := entry.String()
	if err != nil {
		return fmt.Errorf("ERROR_WHILE_GETTING_STRING:%w", err)
	}

	for _, w := range hook.Writer {
		_, err = w.Write([]byte(line))
	}

	return err // nolint
}

func (hook *writerHook) Levels() []logrus.Level {
	return hook.LogLevels
}

var e *logrus.Entry // nolint

type Logger struct {
	*logrus.Entry
}

func Get() *Logger {
	Init()

	return &Logger{e}
}

func Init() {
	log := logrus.New()
	log.SetReportCaller(true)
	log.Formatter = &logrus.TextFormatter{ // nolint
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			filename := path.Base(f.File)

			return fmt.Sprintf("%s:%d", filename, f.Line), fmt.Sprintf("%s()", f.Function)
		},
		DisableColors: false,
		FullTimestamp: true,
	}

	allFile, err := os.OpenFile("./logs/all.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o640)
	if err != nil {
		panic(fmt.Sprintf("[Error]: %s", err))
	}

	log.SetOutput(ioutil.Discard) // Send all logs to nowhere by default

	log.AddHook(&writerHook{
		Writer:    []io.Writer{allFile, os.Stdout},
		LogLevels: logrus.AllLevels,
	})

	log.SetLevel(logrus.TraceLevel)

	e = logrus.NewEntry(log)
}
