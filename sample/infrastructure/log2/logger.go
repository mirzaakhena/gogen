package log2

import (
	"context"
	"fmt"
	"log"
)

// LogPrinter general interface for printing the log
type LogPrinter interface {
	LogPrint(ctx context.Context, flag string, message string)
}

// private variable to store the implementation
var logPrinterInstance LogPrinter

// init instatiate the default implementation
func init() {
	SetLogPrinter(&logPrinterDefault{})
}

// SetLogPrinter changing the log implementation. lg must not nil
func SetLogPrinter(lg LogPrinter) {
	if lg != nil {
		logPrinterInstance = lg
	}
}

// logPrinterDefault is default implementation of LogPrinter
type logPrinterDefault struct {
}

// LogPrint simply print the message to console
func (r *logPrinterDefault) LogPrint(ctx context.Context, flag string, message string) {
	log.Printf("%s %s\n", flag, message)
}

// Info is general info log
func Info(ctx context.Context, message string, args ...interface{}) {
	messageWithArgs := fmt.Sprintf(message, args...)
	logPrinterInstance.LogPrint(ctx, "INF", messageWithArgs)
}

// Error is general error log
func Error(ctx context.Context, message string, args ...interface{}) {
	messageWithArgs := fmt.Sprintf(message, args...)
	logPrinterInstance.LogPrint(ctx, "ERR", messageWithArgs)
}
