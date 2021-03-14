package gr8log

import "time"

type Entry interface {
	Level() Level
	Message() string
	Time() time.Time
}

func NewEntry(level Level, msg string, args ...Args) Entry {
	entry := mapEntry{
		"level":   level,
		"message": msg,
		"time":    time.Now(),
	}

	for _, arg := range args {
		for k, v := range arg {
			if k == "level" || k == "message" || k == "time" {
				continue
			}

			entry[k] = v
		}
	}

	return entry
}

type mapEntry map[string]interface{}

func (m mapEntry) Level() Level {
	return m["level"].(Level)
}

func (m mapEntry) Message() string {
	return m["message"].(string)
}

func (m mapEntry) Time() time.Time {
	return m["time"].(time.Time)
}
