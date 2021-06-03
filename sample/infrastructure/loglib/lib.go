package loglib

import (
	"accounting/infrastructure/log"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"os"
	"time"
)

type dataContextType string

const (
	logGroupIDType dataContextType = "LOG_GROUP_ID"
	startTimeType  dataContextType = "START_TIME"
)

const (
	logGroupIDInitialValue = "00000000000000"
	startTimeInitialValue  = "000000.000"
)

var GenerateLogGroupID func() func() string

func init() {
	GenerateLogGroupID = generateLogGroupIDDefault

	//LogWithPlainFormat()
}

// SetOutput
// f, err := os.OpenFile("hello.txt", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
// defer f.Close()
// if err != nil {
// 	panic(err)
// }
func (r *baseWriteContext) SetOutput(o io.Writer) {
	if o == nil {
		return
	}
	r.out = o
}

type baseWriteContext struct {
	out io.Writer
}

func GetLogWithPlainFormat() *logPrinterPlain {
	x := &logPrinterPlain{}
	x.SetOutput(os.Stdout)
	return x
}

func GetLogWithJSONFormat() *logPrinterJSON {
	x := &logPrinterJSON{}
	x.SetOutput(os.Stdout)
	return x
}

type logPrinterPlain struct {
	baseWriteContext
}

func (r *baseWriteContext) WriteContext(ctx context.Context, data ...interface{}) context.Context {

	ctx = context.WithValue(ctx, startTimeType, getStartTimeFormatted())

	logGroupIDValue := ctx.Value(logGroupIDType)
	if logGroupIDValue == nil && GenerateLogGroupID != nil {
		f := GenerateLogGroupID()
		ctx = context.WithValue(ctx, logGroupIDType, f())
	}

	return ctx
}

func (r *logPrinterPlain) LogPrint(ctx context.Context, flag string, message interface{}) {
	fmt.Fprintf(r.out, "%s %s %s %s %s %s\n", getCurrentTimeFormatted(), extractStartTime(ctx), extractLogGroupID(ctx), flag, log.GetFileLocationInfo(3), message)
}

type logPrinterJSON struct {
	baseWriteContext
}

func (r *logPrinterJSON) LogPrint(ctx context.Context, flag string, message interface{}) {
	info := map[string]interface{}{
		"flag":    flag,
		"message": message,
		"date":    getCurrentTimeFormatted(),
		"start":   extractStartTime(ctx),
		"thisid":  extractLogGroupID(ctx),
		"file":    log.GetFileLocationInfo(3),
	}
	m, _ := json.Marshal(info)
	fmt.Fprintf(r.out, "%v\n", string(m))
}

// getCurrentTimeFormatted get the time formatted
func getCurrentTimeFormatted() string {
	return fmt.Sprintf("%s", time.Now().Format("0102 150405.000"))
}

// getStartTimeFormattedDefault get the time formatted
func getStartTimeFormatted() string {
	return fmt.Sprintf("%s", time.Now().Format("0405.000"))
}

// extractLogGroupIDDefault extract the log group id from context
func extractLogGroupID(ctx context.Context) string {

	if ctx == nil {
		return logGroupIDInitialValue
	}

	logGroupIDInterface := ctx.Value(logGroupIDType)
	if logGroupIDInterface == nil {
		return logGroupIDInitialValue
	}

	logGroupID, _ := logGroupIDInterface.(string)

	return logGroupID
}

// extractStartTime extract the start time from context
func extractStartTime(ctx context.Context) string {

	if ctx == nil {
		return startTimeInitialValue
	}

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
