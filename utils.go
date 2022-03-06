package main

import (
	"context"
	"math/rand"
	"os"
	"reflect"
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
	ctx, cancel = context.WithTimeout(ctx, 40*time.Second)

	return &ctx, cancel
}

func ChunkSlice(slice interface{}, chunkSize int) interface{} {
	sliceType := reflect.TypeOf(slice)
	sliceVal := reflect.ValueOf(slice)
	length := sliceVal.Len()
	if sliceType.Kind() != reflect.Slice {
		panic("parameter must be []T")
	}
	n := 0
	if length%chunkSize > 0 {
		n = 1
	}
	SST := reflect.MakeSlice(reflect.SliceOf(sliceType), 0, length/chunkSize+n)
	st, ed := 0, 0
	for st < length {
		ed = st + chunkSize
		if ed > length {
			ed = length
		}
		SST = reflect.Append(SST, sliceVal.Slice(st, ed))
		st = ed
	}
	return SST.Interface()
}

func ShuffleSlice(slice []string) {
	source := rand.NewSource(time.Now().UnixNano())
	random := rand.New(source)
	for i := len(slice) - 1; i > 0; i-- {
		j := random.Intn(i + 1)
		slice[i], slice[j] = slice[j], slice[i]
	}
}
