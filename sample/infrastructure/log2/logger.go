package log2

import (
	"context"
	"fmt"
	"log"
	"runtime"
	"strings"
)

type traceDataType int

const traceDataKey traceDataType = 1

// LogPrinter general interface for printing the log
type LogPrinter interface {
	LogPrint(ctx context.Context, flag string, message string)
	WriteContext(ctx context.Context, traceData string) context.Context
}

func Context(ctx context.Context, appData string) context.Context {
	return logPrinterInstance.WriteContext(ctx, appData)
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

func (r *logPrinterDefault) WriteContext(ctx context.Context, traceData string) context.Context {
	return context.WithValue(ctx, traceDataKey, traceData)
}

// LogPrint simply print the message to console
func (r *logPrinterDefault) LogPrint(ctx context.Context, flag string, message string) {
	if ctx != nil {
		if v := ctx.Value(traceDataKey); v != nil {
			log.Printf("[%s] %s %s %s\n", flag, v, GetFileLocationInfo(3), message)
			return
		}
	}
	log.Printf("[%s] %s %s\n", flag, GetFileLocationInfo(3), message)
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

// GetFileLocationInfo get the function information like filename and line number
// skip is the parameter that need to adjust if we add new method layer
func GetFileLocationInfo(skip int) string {
	pc, _, line, ok := runtime.Caller(skip)
	if !ok {
		return ""
	}
	funcName := runtime.FuncForPC(pc).Name()
	x := strings.LastIndex(funcName, "/")
	return fmt.Sprintf("%s:%d", funcName[x+1:], line)
}
