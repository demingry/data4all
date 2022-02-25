package main

import (
	"context"
	"fmt"
	"sync"

	"github.com/chromedp/cdproto/dom"
	"github.com/chromedp/chromedp"
)

var (
	threads  chan struct{}
	t        int64
	mu       sync.Mutex
	transfer chan interface{}
)

func sourceFromDriver(url string) {

	ctx, cancel := InitDriver()
	defer cancel()

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
		return
	}

	fmt.Println(string(res))

	<-threads

}

func main() {

	// threads_number, _ := strconv.Atoi(os.Getenv("THREADS_NUMBER"))
	// threads = make(chan struct{}, threads_number)

	// params := strings.Split(os.Getenv("REQUEST_PARAMS"), "-")
	// start, _ := strconv.Atoi(params[0])
	// end, _ := strconv.Atoi(params[1])

	// sourceNodes := make(map[string]interface{})
	// for i := start; i < end; i++ {
	// 	threads <- struct{}{}
	// 	ctx, cancel := InitDriver()
	// 	nodes_instance := NewNodes()
	// 	go nodes_instance.Execute(
	// 		`https://dataverse.harvard.edu/dataverse/harvard?q=&sort=dateSort&order=desc&types=datasets&page=`+fmt.Sprintf("%d", i),
	// 		`.card-title-icon-block a`,
	// 		ctx,
	// 		cancel,
	// 		`href`,
	// 		&sourceNodes,
	// 	)
	// }

	// for {
	// 	if len(threads) == 0 {
	// 		break
	// 	}
	// }

	// var sourceElements []interface{}
	// selectors := make(map[string]string)
	// selectors[`title`] = `span#title`
	// selectors[`description`] = `#dsDescription div`
	// for _, i := range sourceNodes["NodesValue"].([]string) {
	// 	threads <- struct{}{}
	// 	elements_instance := NewElements()
	// 	ctx, cancel := InitDriver()
	// 	go elements_instance.Execute(
	// 		`https://dataverse.harvard.edu`+i,
	// 		selectors,
	// 		ctx,
	// 		cancel,
	// 		&sourceElements,
	// 	)
	// }

	// var sourcePage []string
	// for _, i := range sourceNodes["NodesValue"].([]string) {
	// 	threads <- struct{}{}
	// 	page_instance := NewPage(`PageFromDriver`)
	// 	ctx, cancel := InitDriver()
	// 	go page_instance.Execute(
	// 		`https://dataverse.harvard.edu`+i,
	// 		ctx,
	// 		cancel,
	// 		&sourcePage,
	// 	)

	// 	time.Sleep(3 * time.Second)

	// }

	href := []string{
		`https://dataverse.harvard.edu/dataset.xhtml?persistentId=doi:10.7910/DVN/17XS9I`,
		`https://dataverse.harvard.edu/dataset.xhtml?persistentId=doi:10.7910/DVN/UM5S3X`,
		`https://dataverse.harvard.edu/dataset.xhtml?persistentId=doi:10.7910/DVN/UA8AGD`,
		`https://dataverse.harvard.edu/dataset.xhtml?persistentId=doi:10.7910/DVN/NKCQM1`,
	}
	threads := make(chan struct{}, 3)
	for _, v := range href {
		threads <- struct{}{}
		fmt.Println(v)
		sourceFromDriver(v)
	}

	for {
		if len(threads) == 0 {
			break
		}
	}

}
