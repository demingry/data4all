package data4all

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
		go nodes_instance.Execute(`https://dataverse.harvard.edu/dataverse/harvard?q=&sort=dateSort&order=desc&types=datasets&page=`+fmt.Sprintf("%d", i), `.card-title-icon-block a`, `href`)
	}

	nodesValue := nodes_instance.(IGetter).Getter()[1]
	fmt.Println(nodesValue.([]string))

	fmt.Println("Finished")
	// sitemap_instance := NewSitemap()
	// sitemap_instance.Execute()

	// nodesValue := nodes_instance.(IGetter).Getter()[1]
	// for _, i := range nodesValue.([]string) {
	// 	threads <- struct{}{}
	// 	elements_instance := NewElements()
	// 	go elements_instance.Execute(i)
	// }

	for {
		if len(threads) == 0 {
			break
		}
	}

}
