package logger

import (
	"io"
	"os"
	"sync"
	"time"

	"github.com/rs/zerolog"
	"gopkg.in/natefinch/lumberjack.v2"

	"github.com/VladPetriv/tg_scanner/config"
)

type Logger struct {
	*zerolog.Logger
}

func newFileWriter(filename string) io.Writer {
	return &lumberjack.Logger{
		Filename: filename,
	}
}

func Get(cfg *config.Config) *Logger {
	var (
		logger Logger
		once   sync.Once
	)

	once.Do(func() {
		// By default create console writer
		writers := []io.Writer{zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.Stamp}}

		if cfg.LogFilename != "" {
			writers = append(writers, newFileWriter(cfg.LogFilename))
		}

		if cfg.LogLevel != "" {
			level, err := zerolog.ParseLevel(cfg.LogLevel)
			if err != nil {
				panic(err)
			}

			zerolog.SetGlobalLevel(level)
		}

		multiWriters := io.MultiWriter(writers...)

		zeroLogger := zerolog.New(multiWriters).With().Timestamp().Logger()

		logger = Logger{&zeroLogger}
	})

	return &logger
}
