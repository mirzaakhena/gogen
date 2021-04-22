package log2

import (
	"context"
	"fmt"
)

// LogPrinter general interface for printing the log
type LogPrinter interface {
	LogPrint(ctx context.Context, flag string, data interface{})
	WriteContext(ctx context.Context, data ...interface{}) context.Context
}

// private variable to store the implementation
var logPrinterInstance LogPrinter = &logPrinterDefault{}

// SetLogPrinter changing the log implementation. lg must not nil
func SetLogPrinter(lg LogPrinter) {
	if lg != nil {
		logPrinterInstance = lg
	}
}

// Info is general info log
func Info(ctx context.Context, message string, args ...interface{}) {
	messageWithArgs := fmt.Sprintf(message, args...)
	logPrinterInstance.LogPrint(ctx, "INFO ", messageWithArgs)
}

// Error is general error log
func Error(ctx context.Context, message string, args ...interface{}) {
	messageWithArgs := fmt.Sprintf(message, args...)
	logPrinterInstance.LogPrint(ctx, "ERROR", messageWithArgs)
}

// Context called first time to initiate the log data to be passed to other log
func Context(ctx context.Context, data ...interface{}) context.Context {
	return logPrinterInstance.WriteContext(ctx, data)
}
