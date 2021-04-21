package main

import (
	"accounting/infrastructure/log2"
	"accounting/infrastructure/loglib"
	"context"
	"os"
)

func main() {

	x := loglib.GetLogWithPlainFormat()
	f, _ := os.OpenFile("test.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	defer f.Close()
	x.SetOutput(f)
	log2.SetLogPrinter(x)
	ctx := log2.Context(context.Background(), "createJurnal")
	log2.Info(ctx, "Hello")
}
