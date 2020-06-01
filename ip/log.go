package ip

import "io"

type Logger interface {
	SetOutput(w io.Writer)
	Printf(format string, v ...interface{})
	Print(v ...interface{})
	Println(v ...interface{})
	SetPrefix(prefix string)
}