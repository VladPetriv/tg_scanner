package logger

import (
	"fmt"
	"io"
	"os"
	"path"
	"runtime"

	"github.com/VladPetriv/tg_scanner/pkg/config"
	"github.com/sirupsen/logrus"
)

var e *logrus.Entry // nolint

type Logger struct {
	*logrus.Entry
}

func Get() *Logger {
	cfg, err := config.Get()
	if err != nil {
		panic(err)
	}

	Init(cfg.LogLevel)

	return &Logger{e}
}

func Init(logLevel string) {
	level, err := logrus.ParseLevel(logLevel)
	if err != nil {
		panic(err)
	}

	log := logrus.New()

	log.SetReportCaller(true)

	log.Formatter = &logrus.JSONFormatter{ // nolint
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			filename := path.Base(f.File)

			return fmt.Sprintf("%s:%d", filename, f.Line), fmt.Sprintf("%s()", f.Function)
		},
		PrettyPrint: true,
	}

	file, err := os.OpenFile("./logs/all.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o640)
	if err != nil {
		panic(fmt.Sprintf("failed to open file with logs: %s", err))
	}

	log.SetOutput(io.MultiWriter(file, os.Stdout))

	log.SetLevel(level)

	e = logrus.NewEntry(log)
}
