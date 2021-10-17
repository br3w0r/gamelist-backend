package logger

import (
	"fmt"
	"os"
	"time"
)

type logger struct{}

var Logger = new(logger)

func (l *logger) Write(p []byte) (n int, err error) {
	return fmt.Fprintf(os.Stdout, "(%d) %s", time.Now().Unix(), p)
}
