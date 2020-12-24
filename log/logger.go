package gr8log

import (
	"strings"
)

type Level string

const (
	LevelDebug Level = "debug"
	LevelInfo  Level = "info"
	LevelWarn  Level = "warn"
	LevelError Level = "error"
	LevelFatal Level = "fatal"
)

func LevelFromString(level string) Level {
	switch strings.ToLower(level) {
	case "debug":
		return LevelDebug
	case "info":
		return LevelInfo
	case "warn":
		return LevelWarn
	case "error":
		return LevelError
	case "fatal":
		return LevelFatal
	}

	return LevelDebug
}

// Logger is an interface for leveled logging
type Logger interface {
	Log(Entry)
	WithArgs(Args) Logger
}

type Args map[string]interface{}
