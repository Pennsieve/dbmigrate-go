package dbmigrate

import "log"

// logger implements migrate.Logger; if we don't pass migrate.Migrate one of
// these, it won't do its internal logging.
type logger struct {
	IsVerbose bool
}

func newLogger(verbose bool) *logger {
	return &logger{IsVerbose: verbose}
}

func (l *logger) Printf(format string, v ...interface{}) {
	log.Printf(format, v...)
}

func (l *logger) Verbose() bool {
	return l.IsVerbose
}
