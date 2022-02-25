package main

import (
	"context"
	"os"
	"time"

	"github.com/chromedp/chromedp"
)

func Finished() {
	<-threads
}

func WriteFile(path string, content []byte) error {

	file, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	if _, err := file.WriteString(string(content) + "\n"); err != nil {
		return err
	}

	return nil

}

func InitDriver() (*context.Context, context.CancelFunc) {

	opt := []func(allocator *chromedp.ExecAllocator){
		chromedp.ExecPath(os.Getenv("GOOGLE_CHROME_SHIM")),
		chromedp.Flag("headless", true),
		chromedp.Flag("disable-dev-shm-usage", true),
		chromedp.Flag("diable-extensions", true),
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("--no-sandbox", true),
	}

	allocatorCtx, _ := chromedp.NewExecAllocator(
		context.Background(),
		append(opt, chromedp.DefaultExecAllocatorOptions[:]...)[:]...,
	)

	ctx, cancel := chromedp.NewContext(allocatorCtx)
	ctx, cancel = context.WithTimeout(ctx, 25*time.Second)

	return &ctx, cancel
}
