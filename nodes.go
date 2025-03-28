package main

import (
	"context"
	"fmt"
	"os"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
)

type Nodes struct {
	Nodes      []*cdp.Node
	NodesValue []string
}

/*
	params[0]url, params[1]selector, params[2]*context.Context,
	params[3]context.CancelFunc, params[4]attrname(optional), params[5]source(return)
	return []*cdp.Node or []string(optional attrvalue)
*/
func (n *Nodes) Execute(params ...interface{}) (interface{}, error) {

	defer Finished()
	defer n.Getter(params[5])

	ctx, ok := params[2].(*context.Context)
	cancel, ok := params[3].(context.CancelFunc)

	if !ok {
		return nil, fmt.Errorf("Wrong type in params")
	}

	if len(params) == 5 {

		nodes := n.findNodes(fmt.Sprintf("%v", params[0]), fmt.Sprintf("%v", params[1]), ctx, cancel)
		if nodes == nil {

			//设置重试次数
			mu.Lock()
			retries++
			mu.Unlock()
			return nil, fmt.Errorf("Error in FindNodes")
		}

		copy(n.Nodes, nodes)
		return nodes, nil

	} else if len(params) == 6 {

		nodes := n.findNodes(fmt.Sprintf("%v", params[0]), fmt.Sprintf("%v", params[1]), ctx, cancel)
		if nodes == nil {

			//设置重试次数
			mu.Lock()
			retries++
			mu.Unlock()
			return nil, fmt.Errorf("Error in FindNodes")
		}

		var attrvalue []string

		for _, i := range nodes {
			v := n.getNodeAttr(i, fmt.Sprintf("%v", params[4]))
			if v == "" {
				continue
			}
			attrvalue = append(attrvalue, v)
		}

		return attrvalue, nil

	}
	return nil, fmt.Errorf("Params exceeded")
}

func (n *Nodes) findNodes(url string,
	selector string,
	ctx *context.Context,
	cancel context.CancelFunc,
) []*cdp.Node {

	var nodes []*cdp.Node

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

	fmt.Println("Scraping page: ", url)

	if err := chromedp.Run(*ctx,
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

/*
	map[string]interface: {"Nodes":[]*cdp.Node, "NodesValue":[]string}
*/
func (n *Nodes) Getter(source interface{}) {

	sourceConver, ok := source.(*map[string]interface{})
	if !ok {
		return
	}

	NodesSlice, _ := (*sourceConver)["Nodes"].([]*cdp.Node)
	NodesSlice = append(NodesSlice, n.Nodes...)
	(*sourceConver)["Nodes"] = NodesSlice

	NodesValueSlice, _ := (*sourceConver)["NodesValue"].([]string)
	NodesValueSlice = append(NodesValueSlice, n.NodesValue...)
	(*sourceConver)["NodesValue"] = NodesValueSlice
}
