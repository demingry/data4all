package data4all

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/chromedp/cdproto/dom"
	"github.com/chromedp/chromedp"
)

type PageFromDriver struct {
	pagesource string
}

type PageFromRaw struct {
	pagesource string
}

/*
	params[0]url, params[1]context.Context, params[2]context.CancelFunc
*/
func (pd *PageFromDriver) Execute(params ...interface{}) (interface{}, error) {

	defer Finished()
	if len(params) < 3 {
		return nil, fmt.Errorf("Not enough params")
	}
	ctx, ok := params[1].(*context.Context)
	cancel, ok := params[2].(context.CancelFunc)
	if !ok {
		return nil, fmt.Errorf("Wrong type in params")
	}

	pd.sourceFromDriver(fmt.Sprintf("%v", params[0]), ctx, cancel)
	return nil, nil
}

/*
	params[0]url
*/
func (pr *PageFromRaw) Execute(params ...interface{}) (interface{}, error) {

	defer Finished()
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
		newurl, err := sproxy.Execute(url, ctx)
		if err != nil {
			fmt.Println(err)
			return ""
		}

		url = fmt.Sprintf("%v", newurl)
	}

	var res string
	err := chromedp.Run(*ctx,
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

func (pd *PageFromDriver) Getter() []interface{} {

	var data []interface{}
	data = append(data, pd.pagesource)
	return data
}

func (pr *PageFromRaw) Getter() []interface{} {

	var data []interface{}
	data = append(data, pr.pagesource)
	return data
}
