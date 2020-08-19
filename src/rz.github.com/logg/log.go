// Package logg provides functions for printing to the log file.
package logg

import (
	"bufio"
	"fmt"
	"os"
)

var (
	logFile = "results.log"
	w       *bufio.Writer
)

func init() {
	f, err := os.Create(logFile)
	if err != nil {
		panic(err)
	}
	w = bufio.NewWriter(f)
}

// Println prints to the log file with a newline at the end.
func Println(a ...interface{}) {
	_, err := w.WriteString(fmt.Sprintln(a...))
	if err != nil {
		panic(err)
	}
}

// Printf prints to the log file using a format.
func Printf(format string, a ...interface{}) {
	_, err := w.WriteString(fmt.Sprintf(format, a...))
	if err != nil {
		panic(err)
	}
}

// Flush flushes the contents of the buffer to the file.
func Flush() {
	err := w.Flush()
	if err != nil {
		panic(err)
	}
}
