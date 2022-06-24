package log

import (
	"github.com/rs/zerolog"
	"os"
)

// NewLogger return logger instance
func NewLogger() *zerolog.Logger {
	logger := zerolog.New(os.Stderr).Output(zerolog.ConsoleWriter{Out: os.Stderr}).With().Timestamp().Logger()
	return &logger
}
