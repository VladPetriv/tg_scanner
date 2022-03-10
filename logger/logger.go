package logger

import (
	"fmt"
	"path"
	"runtime"

	"github.com/VladPetriv/tg_scanner/internal/file"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
)

type Logger struct {
	*logrus.Logger
}

func Get() *Logger {
	logger := Logger{logrus.New()}

	err := file.CreateFilesForLogger("logs")
	if err != nil {
		fmt.Println("error is : ", err)
	}
	pathMap := lfshook.PathMap{
		logrus.InfoLevel:  "./logs/info.log",
		logrus.DebugLevel: "./logs/debug.log",
		logrus.ErrorLevel: "./logs/error.log",
		logrus.FatalLevel: "./logs/fatal.log",
		logrus.PanicLevel: "./logs/panic.log",
		logrus.TraceLevel: "./logs/trace.log",
		logrus.WarnLevel:  "./logs/warning.log",
	}
	logger.AddHook(
		lfshook.NewHook(
			pathMap,
			&logrus.TextFormatter{
				CallerPrettyfier: func(f *runtime.Frame) (function string, file string) {
					fileName := path.Base(f.File)
					return fmt.Sprintf("%s():", f.Func.Name()), fmt.Sprintf("%s:%d", fileName, f.Line)
				},
				FullTimestamp: true,
			},
		),
	)

	return &logger
}
