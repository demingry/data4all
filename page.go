package data4all

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/chromedp/cdproto/dom"
	"github.com/chromedp/chromedp"
)

type PageFromDriver struct {
}

type PageFromRaw struct {
}

func (pd *PageFromDriver) Execute(params ...interface{}) (interface{}, error) {

	return nil, nil
}

func (pr *PageFromRaw) Execute(params ...interface{}) (interface{}, error) {

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

func (pd *PageFromDriver) sourceFromDriver(url string) string {

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

	defer cancel()

	if proxy_list := os.Getenv("PROXY_LIST"); proxy_list != "" {
		sproxy := NewProxy()
		newurl, err := sproxy.Execute(url, ctx)
		if err != nil {
			fmt.Println(err)
			return ""
		}

		url = fmt.Sprintf("%v", newurl)
	}

	var res string
	err := chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.ActionFunc(func(ctx context.Context) error {
			node, err := dom.GetDocument().Do(ctx)
			if err != nil {
				return err
			}
			res, err = dom.GetOuterHTML().WithNodeID(node.NodeID).Do(ctx)
			return err
		}),
	)

	if err != nil {
		fmt.Println(err)
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
