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

// logPrinterDefault is default implementation of LogPrinter
type logPrinterDefault struct {
}

// WriteContext passing data to
func (r *logPrinterDefault) WriteContext(ctx context.Context, data ...interface{}) context.Context {
	if len(data) > 0 {
		return context.WithValue(ctx, traceDataKey, data[0])
	}
	return ctx
}

// LogPrint simply print the message to console
func (r *logPrinterDefault) LogPrint(ctx context.Context, flag string, data interface{}) {
	if ctx != nil {
		if v := ctx.Value(traceDataKey); v != nil {
			log.Printf("[%s] %s %s %s\n", flag, v, GetFileLocationInfo(3), data)
			return
		}
	}
	log.Printf("[%s] %s %s\n", flag, GetFileLocationInfo(3), data)
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
