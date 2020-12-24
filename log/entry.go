package gr8log

type Entry interface {
	Level() Level
	Message() string
}

func NewEntry(level Level, msg string, args ...Args) Entry {
	entry := mapEntry{
		"level":   level,
		"message": msg,
	}

	for _, arg := range args {
		for k, v := range arg {
			if k == "level" || k == "message" {
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
