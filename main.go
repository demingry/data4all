package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
)

var (
	threads chan struct{}
	t       int64
	mu      sync.Mutex
)

func main() {

	threads_number, _ := strconv.Atoi(os.Getenv("THREADS_NUMBER"))
	threads = make(chan struct{}, threads_number)

	params := strings.Split(os.Getenv("REQUEST_PARAMS"), "-")
	start, _ := strconv.Atoi(params[0])
	end, _ := strconv.Atoi(params[1])

	sourceNodes := make(map[string]interface{})
	for i := start; i < end; i++ {
		threads <- struct{}{}
		ctx, cancel := InitDriver()
		nodes_instance := NewNodes()
		go nodes_instance.Execute(
			`https://dataverse.harvard.edu/dataverse/harvard?q=&sort=dateSort&order=desc&types=datasets&page=`+fmt.Sprintf("%d", i),
			`.card-title-icon-block a`,
			ctx,
			cancel,
			`href`,
			&sourceNodes,
		)
	}

	for {
		if len(threads) == 0 {
			break
		}
	}

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

	chunked := ChunkSlice(sourceNodes["NodesValue"], 4)

	for _, i := range chunked.([][]string) {
		var sourcePage []string
		var a int = 0
		for _, v := range i {
			a++
			threads <- struct{}{}
			page_instance := NewPage(`PageFromDriver`)
			ctx, cancel := InitDriver()
			go page_instance.Execute(
				`https://dataverse.harvard.edu`+v,
				ctx,
				cancel,
				&sourcePage,
			)
		}

		for {
			if len(threads) == 0 {
				break
			}
		}

		fmt.Println("Length: " + fmt.Sprintf("%d", a))
		fmt.Println("sourcePage: " + fmt.Sprintf("%d", len(sourcePage)))
	}

	for {
		if len(threads) == 0 {
			break
		}
	}

}
