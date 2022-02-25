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

	nodes_instance := NewNodes()
	for i := start; i < end; i++ {
		threads <- struct{}{}
		ctx, cancel := InitDriver()
		go nodes_instance.Execute(
			`https://dataverse.harvard.edu/dataverse/harvard?q=&sort=dateSort&order=desc&types=datasets&page=`+fmt.Sprintf("%d", i),
			`.card-title-icon-block a`,
			ctx,
			cancel,
			`href`)
	}

	for {
		if len(threads) == 0 {
			break
		}
	}

	selectors := make(map[string]string)
	selectors[`title`] = `span#title`
	nodesValue := nodes_instance.(IGetter).Getter()[1]
	for _, i := range nodesValue.([]string) {
		threads <- struct{}{}
		elements_instance := NewElements()
		ctx, cancel := InitDriver()
		go elements_instance.Execute(
			`https://dataverse.harvard.edu`+i,
			selectors,
			ctx,
			cancel,
		)

		elements := elements_instance.(IGetter).Getter()[0]

		for k, v := range elements.(map[string]string) {
			fmt.Println(k, ":", v)
		}
	}

	for {
		if len(threads) == 0 {
			break
		}
	}

}
