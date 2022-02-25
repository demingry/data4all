package data4all

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
)

type Nodes struct {
	Nodes      []*cdp.Node
	NodesValue []string
}

/*
	params[0]url, params[1]selector, params[2]attrname(optional)
	return []*cdp.Node or []string(optional attrvalue)
*/
func (n *Nodes) Execute(params ...interface{}) (interface{}, error) {

	defer Finished()

	if len(params) == 2 {

		nodes := n.findNodes(fmt.Sprintf("%v", params[0]), fmt.Sprintf("%v", params[1]))
		if nodes == nil {
			return nil, fmt.Errorf("Error in FindNodes")
		}

		copy(n.Nodes, nodes)
		return nodes, nil

	} else if len(params) == 3 {

		nodes := n.findNodes(fmt.Sprintf("%v", params[0]), fmt.Sprintf("%v", params[1]))
		if nodes == nil {
			return nil, fmt.Errorf("Error in FindNodes")
		}

		var attrvalue []string

		for _, i := range nodes {
			v := n.getNodeAttr(i, fmt.Sprintf("%v", params[2]))
			if v == "" {
				continue
			}
			attrvalue = append(attrvalue, v)
		}

		return attrvalue, nil

	}
	return nil, fmt.Errorf("Params exceeded")
}

func (n *Nodes) findNodes(url string, selector string) []*cdp.Node {

	var nodes []*cdp.Node

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
			return nil
		}

		url = fmt.Sprintf("%v", newurl)
	}

	if err := chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.Nodes(selector, &nodes),
	); err != nil {
		fmt.Printf("[!]Err gethref in: %s\n", err.Error())
		return nil
	}

	return nodes

}

func (n *Nodes) getNodeAttr(node *cdp.Node, attrname string) string {

	n.NodesValue = append(n.NodesValue, node.AttributeValue(attrname))
	return node.AttributeValue(attrname)
}

func NewNodes() Icommand {
	return &Nodes{}
}

func (n *Nodes) Getter() []interface{} {
	var data []interface{}
	data = append(data, n.Nodes)
	data = append(data, n.NodesValue)
	return data
}
