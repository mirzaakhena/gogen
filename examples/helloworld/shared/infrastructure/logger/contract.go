package logger

import (
	"context"
	"encoding/json"
	"fmt"
	"runtime"
	"strings"
)

type Logger interface {
	Info(ctx context.Context, message string, args ...any)
	Error(ctx context.Context, message string, args ...any)
}

type traceDataType int

const traceDataKey traceDataType = 1

func SetTraceID(ctx context.Context, traceID string) context.Context {
	return context.WithValue(ctx, traceDataKey, traceID)
}

func GetTraceID(ctx context.Context) string {

	// default traceID
	traceID := "0000000000000000"

	if ctx != nil {
		if v := ctx.Value(traceDataKey); v != nil {
			traceID = v.(string)
		}
	}

	return traceID
}

// getFileLocationInfo get the function information like filename and line number
// skip is the parameter that need to adjust if we add new method layer
func getFileLocationInfo(skip int) string {
	pc, _, line, ok := runtime.Caller(skip)
	if !ok {
		return ""
	}
	funcName := runtime.FuncForPC(pc).Name()
	x := strings.LastIndex(funcName, "/")
	return fmt.Sprintf("%s:%d", funcName[x+1:], line)
}

func toJsonString(obj any) string {
	bytes, _ := json.Marshal(obj)
	return string(bytes)
}
