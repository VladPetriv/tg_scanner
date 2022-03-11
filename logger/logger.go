package logger

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"runtime"

	"github.com/VladPetriv/tg_scanner/internal/file"
	"github.com/sirupsen/logrus"
)

type writerHook struct {
	Writer    []io.Writer
	LogLevels []logrus.Level
}

func (hook *writerHook) Fire(entry *logrus.Entry) error {
	line, err := entry.String()
	if err != nil {
		return err
	}
	for _, w := range hook.Writer {
		_, err = w.Write([]byte(line))
	}
	return err
}

// Levels define on which log levels this hook would trigger
func (hook *writerHook) Levels() []logrus.Level {
	return hook.LogLevels
}

var e *logrus.Entry

type Logger struct {
	*logrus.Entry
}

func Get() *Logger {
	Init()
	return &Logger{e}
}

func Init() {
	l := logrus.New()
	l.SetReportCaller(true)
	l.Formatter = &logrus.TextFormatter{
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			filename := path.Base(f.File)
			return fmt.Sprintf("%s:%d", filename, f.Line), fmt.Sprintf("%s()", f.Function)
		},
		DisableColors: false,
		FullTimestamp: true,
		ForceColors:   true,
	}

	allFile, err := file.CreateFileForLogger("./logs/all.log")
	if err != nil {
		panic(fmt.Sprintf("[Error]: %s", err))
	}

	l.SetOutput(ioutil.Discard) // Send all logs to nowhere by default

	l.AddHook(&writerHook{
		Writer:    []io.Writer{allFile, os.Stdout},
		LogLevels: logrus.AllLevels,
	})

	l.SetLevel(logrus.TraceLevel)

	e = logrus.NewEntry(l)
}
