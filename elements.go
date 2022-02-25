package data4all

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/chromedp/chromedp"
)

type Elemtns struct {
	elements map[string]string
}

/*
	params[0]url, params[1]selectors
*/
func (e *Elemtns) Execute(params ...interface{}) (interface{}, error) {

	selectors, ok := params[1].(map[string]string)

	if !ok {
		return nil, fmt.Errorf("Wrong map[string]string type in params")
	}
	elements := e.findElements(fmt.Sprintf("%v", params[0]), selectors)
	return elements, nil
}

func (e *Elemtns) findElements(url string, selectors map[string]string) map[string]string {

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
	ctx, cancel = context.WithTimeout(ctx, 100*time.Second)

	defer cancel()
	if proxy_list := os.Getenv("PROXY_LIST"); proxy_list != "" {

		sproxy := NewProxy()
		newurl, err := sproxy.Execute(url, ctx)

		if err != nil {
			fmt.Println(err)
			return nil
		}

		url = fmt.Sprintf("%v", newurl)
	}

	var foundElemets map[string]string
	if err := chromedp.Run(ctx,
		chromedp.Navigate(url),
	); err != nil {
		fmt.Println("Navigate: ", err)
	}

	for k, v := range selectors {

		var tmp string
		ctxchild, _ := context.WithTimeout(ctx, 15*time.Second)
		if err := chromedp.Run(ctxchild,
			chromedp.Text(v, &tmp, chromedp.BySearch),
		); err != nil {
			fmt.Println(fmt.Sprintf("%s", k), err)
		}
		foundElemets[k] = tmp
	}

	e.checkRes(foundElemets)

	return foundElemets
}

func NewElements() Icommand {
	return &Elemtns{}
}

func (e *Elemtns) checkRes(res map[string]string) {

}

func (e *Elemtns) Getter() []interface{} {

	var data []interface{}
	data = append(data, e.elements)
	return data
}
