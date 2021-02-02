package log

import (
	"context"
	"fmt"
	"log"
	"runtime"
	"strings"
)

func Info(ctx context.Context, message string, args ...interface{}) {
	additionalFunctionName := fmt.Sprintf("[INFO] %s ", getFunctionCall(2))
	log.Printf(additionalFunctionName+message, args...)
}

func Error(ctx context.Context, message string, args ...interface{}) {
	additionalFunctionName := fmt.Sprintf("[ERRO] %s ", getFunctionCall(2))
	log.Printf(additionalFunctionName+message, args...)
}

func getFunctionCall(functionLevel int) string {
	pc, _, line, ok := runtime.Caller(functionLevel)
	if !ok {
		return ""
	}
	funcName := runtime.FuncForPC(pc).Name()
	x := strings.LastIndex(funcName, "/")
	return fmt.Sprintf("%s:%d", funcName[x+1:], line)
}
