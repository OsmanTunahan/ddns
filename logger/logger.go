package logger

import (
	"log"
	"os"
)

type Logger struct {
	*log.Logger
}

func NewLogger(level string) *Logger {
	return &Logger{
		Logger: log.New(os.Stdout, "["+level+"] ", log.LstdFlags),
	}
}
