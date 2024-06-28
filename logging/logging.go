package logging

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var Logger zerolog.Logger

func init() {
	logFile := setupLogFile()
	writer := zerolog.MultiLevelWriter(zerolog.ConsoleWriter{Out: os.Stdout}, logFile)

	zerolog.TimeFieldFormat = time.RFC3339
	// zerolog.SetGlobalLevel(zerolog.InfoLevel)
	zerolog.SetGlobalLevel(zerolog.DebugLevel)

	Logger = zerolog.New(writer).With().Timestamp().Caller().Logger()

	log.Logger = Logger
}

func setupLogFile() *os.File {
	logDir := filepath.Join(os.TempDir(), "tabload_logs")

	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		if err := os.MkdirAll(logDir, 0755); err != nil {
			panic(fmt.Sprintf("failed to create log directory: %v", err))
		}
	}

	if _, err := os.Stat(filepath.Join(logDir, "tabload.log")); os.IsNotExist(err) {
		if _, err := os.Create(filepath.Join(logDir, "tabload.log")); err != nil {
			panic(fmt.Sprintf("failed to create log file: %v", err))
		}
	}

	logPath := filepath.Join(logDir, "tabload.log")
	file, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(fmt.Sprintf("failed to open log file: %v", err))
	}

	return file
}

func Debug(msg string) {
	_, file, line, _ := runtime.Caller(1)
	Logger.Debug().Str("file", filepath.Base(file)).Int("line", line).Msg(msg)
}

func Info(msg string) {
	_, file, line, _ := runtime.Caller(1)
	Logger.Info().Str("file", filepath.Base(file)).Int("line", line).Msg(msg)
}

func Warn(msg string) {
	_, file, line, _ := runtime.Caller(1)
	Logger.Warn().Str("file", filepath.Base(file)).Int("line", line).Msg(msg)
}

func Error(msg string, err error) {
	_, file, line, _ := runtime.Caller(1)
	Logger.Error().Err(err).Str("file", filepath.Base(file)).Int("line", line).Msg(msg)
}

func Fatal(msg string, err error) {
	_, file, line, _ := runtime.Caller(1)
	Logger.Fatal().Err(err).Str("file", filepath.Base(file)).Int("line", line).Msg(msg)
}

// a function which returns a stream of logs that the UI can display
func StreamLogs() chan string {
	logStream := make(chan string)

	go func() {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
		defer func() {
			log.Logger = Logger
			close(logStream)
		}()

		log.Logger = log.Output(zerolog.MultiLevelWriter(zerolog.ConsoleWriter{Out: os.Stdout}, os.Stdout))

		// If new logs are created, send them to the logStream
		for {
			select {
			case <-time.After(1 * time.Second):
				continue
			}
		}
	}()

	return logStream
}
