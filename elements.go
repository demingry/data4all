package main

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
	params[0]url, params[1]selectors(map[string]string), params[2]context.Context,
	params[3]context.CancelFunc, params[4]source(return)
*/
func (e *Elemtns) Execute(params ...interface{}) (interface{}, error) {

	defer Finished()
	defer e.Getter(params[4])
	selectors, ok := params[1].(map[string]string)
	ctx, ok := params[2].(*context.Context)
	cancel, ok := params[3].(context.CancelFunc)
	if !ok {
		return nil, fmt.Errorf("Wrong type in params")
	}
	elements := e.findElements(fmt.Sprintf("%v", params[0]), selectors, ctx, cancel)
	return elements, nil
}

func (e *Elemtns) findElements(url string,
	selectors map[string]string,
	ctx *context.Context,
	cancel context.CancelFunc,
) map[string]string {

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

	foundElemets := make(map[string]string)
	if err := chromedp.Run(*ctx,
		chromedp.Navigate(url),
	); err != nil {
		fmt.Println("Navigate: ", err)
	}

	for k, v := range selectors {

		var tmp string
		ctxchild, _ := context.WithTimeout(*ctx, 15*time.Second)
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

func (e *Elemtns) clickElements(selector string, ctx *context.Context) {

	if err := chromedp.Run(*ctx,
		chromedp.Click(selector, chromedp.BySearch),
	); err != nil {
		fmt.Println("Click: ", err)
		return
	}
}

func NewElements() Icommand {
	return &Elemtns{}
}

func (e *Elemtns) checkRes(res map[string]string) {

}

func (e *Elemtns) Getter(source interface{}) {

	sourceConver, ok := source.(*map[string]string)
	if !ok {
		fmt.Println("this1")
		return
	}

	for k, v := range e.elements {
		(*sourceConver)[k] = v
	}
}
