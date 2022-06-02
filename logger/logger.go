package logger

import (
	"fmt"
	"io"
	"os"
	"path"
	"runtime"

	"github.com/sirupsen/logrus"
)

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

	log.SetOutput(io.MultiWriter(allFile, os.Stdout))

	log.SetLevel(logrus.TraceLevel)

	e = logrus.NewEntry(log)
}
