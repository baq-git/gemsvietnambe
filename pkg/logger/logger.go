package logger

import (
	"fmt"
	"os"

	"github.com/rs/zerolog"
)

type Logger struct {
	logger zerolog.Logger
}

var Log Logger

func Init(level, output string) error {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	switch level {
	case "info":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case "warn":
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case "error":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	case "debug":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case "fatal":
		zerolog.SetGlobalLevel(zerolog.FatalLevel)
	case "panic":
		zerolog.SetGlobalLevel(zerolog.PanicLevel)
	}

	if output == "stdout" && os.Getenv("ENV") != "PROD" {
		Log.logger = zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr, NoColor: false}).With().Timestamp().Logger()
	} else {
		if output == "stdout" {
			Log.logger = zerolog.New(os.Stdout).With().Timestamp().Logger()
		} else if output == "stderr" {
			Log.logger = zerolog.New(os.Stderr).With().Timestamp().Logger()
		} else {
			f, err := os.OpenFile(output, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				fmt.Println(err)
				return err
			}
			Log.logger = zerolog.New(f).With().Timestamp().Logger()
		}
	}

	return nil
}

func Info(msg string) {
	Log.logger.Info().Msg(msg)
}

func Warn(msg string) {
	Log.logger.Warn().Msg(msg)
}

func Error(msg string, err error) {
	Log.logger.Err(err).Msg(msg)
}

func Fatal(msg string, err error) {
	Log.logger.Fatal().Err(err).Msg(msg)
}

func Debug(msg string) {
	Log.logger.Debug().Msg(msg)
}

func Panic(msg string, err error) {
	Log.logger.Panic().Err(err).Msg(msg)
}
