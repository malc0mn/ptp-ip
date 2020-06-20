package ip

import (
	"errors"
	"io"
	"log"
)

const (
	// LevelSilent disables all log messages.
	LevelSilent LogLevel = 0
	// LevelVerbose will output warning and error messages.
	LevelVerbose LogLevel = 1
	// LevelVeryVerbose will output info, warning and error messages.
	LevelVeryVerbose LogLevel = 2
	// LevelVeryVerbose will output debug, info, warning and error messages.
	LevelDebug LogLevel = 3
)

type LogLevel byte

// Set() implements flags.Value interface.
func (l *LogLevel) Set(s string) error {
	*l = LevelSilent
	switch s {
	case "v":
		*l = LevelVerbose
	case "vv":
		*l = LevelVeryVerbose
	case "vvv":
		*l = LevelDebug
	default:
		return errors.New("unknown log level")
	}

	return nil
}

// String() implement flags.Value interface.
func (l *LogLevel) String() string {
	switch *l {
	case LevelVerbose:
		return "v"
	case LevelVeryVerbose:
		return "vv"
	case LevelDebug:
		return "vvv"
	}

	return ""
}

// Logger is the interface allowing you to create a custom logger.
type Logger interface {
	Debug(v ...interface{})
	Debugf(format string, v ...interface{})
	Debugln(v ...interface{})
	Error(v ...interface{})
	Errorf(format string, v ...interface{})
	Errorln(v ...interface{})
	Fatal(v ...interface{})
	Fatalf(format string, v ...interface{})
	Fatalln(v ...interface{})
	Info(v ...interface{})
	Infof(format string, v ...interface{})
	Infoln(v ...interface{})
	Warn(v ...interface{})
	Warnf(format string, v ...interface{})
	Warnln(v ...interface{})
}

// StdLogger is the standard logger, a wrapper around the golang log package.
type StdLogger struct {
	level LogLevel
	*log.Logger
}

func (sl *StdLogger) Debug(v ...interface{}) {
	if sl.level >= LevelDebug {
		log.Print(v...)
	}
}

func (sl *StdLogger) Debugf(format string, v ...interface{}) {
	if sl.level >= LevelDebug {
		log.Printf(format, v...)
	}
}

func (sl *StdLogger) Debugln(v ...interface{}) {
	if sl.level >= LevelDebug {
		log.Println(v...)
	}
}

func (sl *StdLogger) Error(v ...interface{}) {
	if sl.level > LevelSilent {
		log.Print(v...)
	}
}

func (sl *StdLogger) Errorf(format string, v ...interface{}) {
	if sl.level > LevelSilent {
		log.Printf(format, v...)
	}
}

func (sl *StdLogger) Errorln(v ...interface{}) {
	if sl.level > LevelSilent {
		log.Println(v...)
	}
}

func (sl *StdLogger) Info(v ...interface{}) {
	if sl.level >= LevelVeryVerbose {
		log.Print(v...)
	}
}

func (sl *StdLogger) Infof(format string, v ...interface{}) {
	if sl.level >= LevelVeryVerbose {
		log.Printf(format, v...)
	}
}

func (sl *StdLogger) Infoln(v ...interface{}) {
	if sl.level >= LevelVeryVerbose {
		log.Println(v...)
	}
}

func (sl *StdLogger) Warn(v ...interface{}) {
	if sl.level >= LevelVerbose {
		log.Print(v...)
	}
}

func (sl *StdLogger) Warnf(format string, v ...interface{}) {
	if sl.level >= LevelVerbose {
		log.Printf(format, v...)
	}
}

func (sl *StdLogger) Warnln(v ...interface{}) {
	if sl.level >= LevelVerbose {
		log.Println(v...)
	}
}

// NewLogger creates a new StdLogger. The out variable sets the destination to which log data will be written.
// The level determines the type of log messages being output.
// The prefix appears at the beginning of each generated log line, or after the log header if the log.Lmsgprefix flag is
// provided.
// The flag argument defines the logging properties.
func NewLogger(level LogLevel, out io.Writer, prefix string, flag int) Logger {
	return &StdLogger{
		level:  level,
		Logger: log.New(out, prefix, flag),
	}
}
