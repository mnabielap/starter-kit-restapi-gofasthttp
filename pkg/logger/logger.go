package logger

import (
	"log"
	"os"
)

// Info logs informational messages
func Info(format string, v ...interface{}) {
	log.Printf("[INFO] "+format, v...)
}

// Error logs error messages
func Error(format string, v ...interface{}) {
	log.Printf("[ERROR] "+format, v...)
}

// Fatal logs error messages and exits
func Fatal(format string, v ...interface{}) {
	log.Printf("[FATAL] "+format, v...)
	os.Exit(1)
}