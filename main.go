package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
)

var (
	threads chan struct{}
	mu      sync.Mutex
)

func main() {

	f, err := os.Create(`results`)
	if err != nil {
		log.Fatal(err)
	}
	f.Close()

	threads_number, _ := strconv.Atoi(os.Getenv("THREADS_NUMBER"))
	threads = make(chan struct{}, threads_number)

	params := strings.Split(os.Getenv("REQUEST_PARAMS"), "-")
	start, _ := strconv.Atoi(params[0])
	end, _ := strconv.Atoi(params[1])

	sourceNodes := make(map[string]interface{})
	for i := start; i < end; i++ {
		for q := 1; q < 40; q++ {

			fmt.Println("Scraped Page: " + fmt.Sprintf("%d", q) + " year: " + fmt.Sprintf("%d", i))
			threads <- struct{}{}
			ctx, cancel := InitDriver()
			nodes_instance := NewNodes()
			go nodes_instance.Execute(
				`https://search.datacite.org/repositories?year=`+fmt.Sprintf("%d", i)+`&page=`+fmt.Sprintf("%d", q),
				`.work a`,
				ctx,
				cancel,
				`href`,
				&sourceNodes,
			)
		}
	}

	for {
		if len(threads) == 0 {
			break
		}
	}

	for k, v := range sourceNodes[`NodesValue`].([]string) {
		sourceNodes[`NodesValue`].([]string)[k] = `https://search.datacite.org/` + v
		WriteFile(`results`, []byte(sourceNodes[`NodesValue`].([]string)[k]))
	}

	upload_instance := NewUpload()
	upload_instance.Execute(`results`, os.Getenv(`INFO`))

}
