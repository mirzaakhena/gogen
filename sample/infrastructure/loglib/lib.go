package loglib

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"os"
	"runtime"
	"strings"
	"time"
)

type contextType string

const (
	logGroupIDType contextType = "LOG_GROUP_ID"
	appNameType    contextType = "APP_NAME"
	startTimeType  contextType = "START_TIME"
)

const (
	logGroupIDInitialValue = "00000000000000"
	appNameInitialValue    = ""
	startTimeInitialValue  = "000000.000"
)

var GenerateLogGroupID func() func() string

var out io.Writer

func init() {
	GenerateLogGroupID = generateLogGroupIDDefault
	SetOutput(os.Stdout)
	//LogWithPlainFormat()
}

// SetOutput
// f, err := os.OpenFile("hello.txt", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
// defer f.Close()
// if err != nil {
// 	panic(err)
// }
func SetOutput(o io.Writer) {
	if o == nil {
		return
	}
	out = o
}

//func LogWithPlainFormat() {
//	log2.SetLogPrinter(&logPrinterPlain{})
//}
//
//func LogWithJSONFormat() {
//	log2.SetLogPrinter(&logPrinterJSON{})
//}

type logPrinterPlain struct{}

func (r *logPrinterPlain) LogPrint(ctx context.Context, flag string, message string) {
	fmt.Fprintf(out, "%s %s %s %s %s %s %s\n", extractStartTime(ctx), getCurrentTimeFormatted(), extractAppName(ctx), extractLogGroupID(ctx), flag, getFileLocationInfo(3), message)
}

type logPrinterJSON struct{}

func (r *logPrinterJSON) LogPrint(ctx context.Context, flag string, message string) {
	info := map[string]interface{}{
		"flag":    flag,
		"message": message,
		"date":    getCurrentTimeFormatted(),
		"start":   extractStartTime(ctx),
		"app":     extractAppName(ctx),
		"id":      extractLogGroupID(ctx),
		"file":    getFileLocationInfo(3),
	}
	m, _ := json.Marshal(info)
	fmt.Fprintf(out, "%v\n", string(m))
}

// getFunctionCall get the function information like filename and line number
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

// getCurrentTimeFormatted get the time formatted
func getCurrentTimeFormatted() string {
	return fmt.Sprintf("%s", time.Now().Format("0102 150405.000"))
}

// getStartTimeFormattedDefault get the time formatted
func getStartTimeFormatted() string {
	return fmt.Sprintf("%s", time.Now().Format("150405.000"))
}

// extractLogGroupIDDefault extract the log group id from context
func extractLogGroupID(ctx context.Context) string {

	logGroupIDInterface := ctx.Value(logGroupIDType)
	if logGroupIDInterface == nil {
		return logGroupIDInitialValue
	}

	logGroupID, _ := logGroupIDInterface.(string)

	return logGroupID
}

// extractAppName extract the app name from context
func extractAppName(ctx context.Context) string {

	appNameValue := ctx.Value(appNameType)
	if appNameValue == nil {
		return appNameInitialValue
	}

	appNameString, _ := appNameValue.(string)

	return appNameString
}

// extractStartTime extract the start time from context
func extractStartTime(ctx context.Context) string {

	startTimeValue := ctx.Value(startTimeType)
	if startTimeValue == nil {
		return startTimeInitialValue
	}
	startTimeStr, _ := startTimeValue.(string)
	return startTimeStr
}

func generateLogGroupIDDefault() func() string {

	const min = 0
	const max = 99999999

	now := time.Now()
	rand.Seed(now.UnixNano())

	return func() string {
		return fmt.Sprintf("%s%08d", now.Format("150405"), rand.Intn(max-min+1)+min)
	}
}

// Context is always called in the beginning request of inport client
// This method will generate the log group id that will distributed to the next method call via context
func Context(ctx context.Context, appName string) context.Context {

	ctx = context.WithValue(ctx, appNameType, appName)

	ctx = context.WithValue(ctx, startTimeType, getStartTimeFormatted())

	logGroupIDValue := ctx.Value(logGroupIDType)
	if logGroupIDValue == nil && GenerateLogGroupID != nil {
		f := GenerateLogGroupID()
		ctx = context.WithValue(ctx, logGroupIDType, f())
	}

	return ctx
}
