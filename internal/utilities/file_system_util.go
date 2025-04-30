package utilities

import (
	"log"
	"runtime"
)

func LogWithLine(path string, message string) {
	// Get the caller information
	_, file, line, ok := runtime.Caller(1)

	if ok {
		log.Printf("%s:%d %s\n", file, line, message)
	} else {
		log.Printf("[%s] %s\n", path, message)
	}
}
