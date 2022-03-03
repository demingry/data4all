package main

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
)

type proxy struct {
}

/*
	params[0]url, params[1]*prevcontext
	return refactor url
*/
func (p *proxy) Execute(params ...interface{}) (interface{}, error) {

	ctx, ok := params[1].(*context.Context)
	if !ok {
		return nil, fmt.Errorf("Wrong context type in params")
	}
	newurl, err := p.doProxy(fmt.Sprintf("%v", params[0]), ctx)
	if err != nil {
		return nil, err
	}

	return newurl, nil
}

func (p *proxy) doProxy(url string, ctx *context.Context) (string, error) {

	if found := strings.Contains(url, `data4all-proxy`); found {
		return url, nil
	}
	proxy_list := strings.Split(os.Getenv("PROXY_LIST"), ";")

	rand.Seed(time.Now().UnixNano())
	randIdx := rand.Intn(len(proxy_list))
	proxy := proxy_list[randIdx] + `.herokuapp.com`

	if err := chromedp.Run(*ctx,
		chromedp.Navigate(`https://`+proxy),
	); err != nil {
		return "", fmt.Errorf("Cannot navigate to proxy")
	}

	url = `https://` + proxy + `/main/` + url
	return url, nil
}

func NewProxy() Icommand {
	return &proxy{}
}
