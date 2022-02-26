package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/chromedp/chromedp"
)

type PageFromDriver struct {
	pagesource string
}

type PageFromRaw struct {
	pagesource string
}

/*
	params[0]url, params[1]context.Context, params[2]context.CancelFunc, params[3]source(return)
*/
func (pd *PageFromDriver) Execute(params ...interface{}) (interface{}, error) {

	defer Finished()
	defer pd.Getter(params[3])
	if len(params) < 4 {
		return nil, fmt.Errorf("Not enough params")
	}
	ctx, ok := params[1].(*context.Context)
	if !ok {
		return nil, fmt.Errorf("Wrong type in params")
	}
	cancel, ok := params[2].(context.CancelFunc)
	if !ok {
		return nil, fmt.Errorf("Wrong type in params")
	}

	pd.pagesource = pd.sourceFromDriver(fmt.Sprintf("%v", params[0]), ctx, cancel)
	return nil, nil
}

/*
	params[0]url, params[1]source(return)
*/
func (pr *PageFromRaw) Execute(params ...interface{}) (interface{}, error) {

	defer Finished()
	defer pr.Getter(params[1])
	pr.sourceFromRaw(fmt.Sprintf("%v", params[0]))
	return nil, nil
}

func (pr *PageFromRaw) sourceFromRaw(url string) string {

	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)")

	resp, err := client.Do(req)

	if err != nil {
		return ""
	}

	if resp.StatusCode != 200 {
		fmt.Println("Http status code:", resp.StatusCode)
		return ""
	}

	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	return string(body)
}

func (pd *PageFromDriver) sourceFromDriver(url string,
	ctx *context.Context,
	cancel context.CancelFunc,
) string {

	defer cancel()

	if proxy_list := os.Getenv("PROXY_LIST"); proxy_list != "" {
		sproxy := NewProxy()
		newurl, err := sproxy.Execute(url, *ctx)
		if err != nil {
			fmt.Println(err)
			return ""
		}

		url = fmt.Sprintf("%v", newurl)
	}

	var res string
	err := chromedp.Run(*ctx,
		chromedp.Navigate(url),
		chromedp.WaitReady("body"),
		chromedp.OuterHTML("html", &res, chromedp.ByQuery),
	)

	if err != nil {
		fmt.Println("Error in get page source: ", err)
		return ""
	}

	return res

}

func NewPage(pagetype string) Icommand {

	if pagetype == "PageFromDriver" {
		return &PageFromDriver{}
	}

	return &PageFromRaw{}
}

/*
	[]string
*/
func (pd *PageFromDriver) Getter(source interface{}) {

	sourceConver, ok := source.(*[]string)
	if !ok {
		return
	}

	*sourceConver = append(*sourceConver, pd.pagesource)

}

/*
	[]string
*/
func (pr *PageFromRaw) Getter(source interface{}) {

	mu.Lock()
	sourceConver, ok := source.(*[]string)
	if !ok {
		return
	}

	*sourceConver = append(*sourceConver, pr.pagesource)
	mu.Unlock()
}
