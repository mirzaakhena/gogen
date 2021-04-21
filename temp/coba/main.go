package main

import (
	"context"
	"github.com/mirzaakhena/gogen2/temp/log/log2"
	"github.com/mirzaakhena/gogen2/temp/log/loglib"
	"os"
	"sync"
)

var wg sync.WaitGroup

func main() {

	loglib.LogWithJSONFormat()

	f, err := os.OpenFile("hello.txt", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	defer f.Close()
	if err != nil {
		panic(err)
	}

	loglib.SetOutput(f)

	{
		ctx := loglib.Context(context.Background(), "satu")

		for i := 0; i < 10; i++ {
			wg.Add(1)
			go sebuahInfo(ctx, i)
		}

	}

	{
		ctx := loglib.Context(context.Background(), "duaa")

		for i := 0; i < 10; i++ {
			wg.Add(1)
			go fungsiLain(ctx, i)
		}

	}

	wg.Wait()

}

func sebuahInfo(ctx context.Context, i int) {
	log2.Info(ctx, "Hello %v", i)
	wg.Done()
}

func fungsiLain(ctx context.Context, i int) {
	log2.Error(ctx, "Lain lagi %v", i)
	wg.Done()
}
