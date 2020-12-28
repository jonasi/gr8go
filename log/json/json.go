package json

import (
	"encoding/json"
	"io"
	"os"
	"sync"

	gr8log "github.com/jonasi/gr8go/log"
	"github.com/mitchellh/mapstructure"
)

type JSONLogger struct {
	Out io.Writer

	initOnce sync.Once
	jsenc    *json.Encoder
	args     []gr8log.Args
}

func (l *JSONLogger) init() {
	l.initOnce.Do(func() {
		if l.Out == nil {
			l.Out = os.Stdout
		}

		l.jsenc = json.NewEncoder(l.Out)
	})
}

func (l *JSONLogger) Log(entry gr8log.Entry) {
	l.init()

	if len(l.args) == 0 {
		l.jsenc.Encode(entry)
		return
	}

	var m map[string]interface{}
	mapstructure.Decode(entry, &m)

	for _, args := range l.args {
		for k, v := range args {
			m[k] = v
		}
	}

	l.jsenc.Encode(m)
}

func (l *JSONLogger) WithArgs(args gr8log.Args) gr8log.Logger {
	return &JSONLogger{Out: l.Out, args: append(l.args, args)}
}

var _ gr8log.Logger = &JSONLogger{}
