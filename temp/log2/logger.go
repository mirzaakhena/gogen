package log2

import (
	"context"
	"fmt"
	"log"
)

type LogPrinter interface {
	LogPrintFormat(ctx context.Context, flag string, message string)
}

var LogPrinterInstance LogPrinter

func init() {
	LogPrinterInstance = &logPrinterDefault{}
}

type logPrinterDefault struct {
}

func (r *logPrinterDefault) LogPrintFormat(ctx context.Context, flag string, message string) {
	log.Printf("%s %s\n", flag, message)
}

// Info is general info log
func Info(ctx context.Context, message string, args ...interface{}) {
	messageWithArgs := fmt.Sprintf(message, args...)
	LogPrinterInstance.LogPrintFormat(ctx, "INFO", messageWithArgs)
}

// Error is general error log
func Error(ctx context.Context, message string, args ...interface{}) {
	messageWithArgs := fmt.Sprintf(message, args...)
	LogPrinterInstance.LogPrintFormat(ctx, "ERRO", messageWithArgs)
}
