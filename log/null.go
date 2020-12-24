package gr8log

var NullLogger = &nullLogger{}

type nullLogger struct{}

var _ Logger = &nullLogger{}

func (l *nullLogger) Log(Entry)            {}
func (l *nullLogger) WithArgs(Args) Logger { return l }
