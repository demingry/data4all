package main

import (
	"encoding/json"
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
		fmt.Println("Scraped Page: " + fmt.Sprintf("%d", i))
		threads <- struct{}{}
		ctx, cancel := InitDriver()
		nodes_instance := NewNodes()
		go nodes_instance.Execute(
			`https://datacatalog.worldbank.org/search?q=&start=`+fmt.Sprintf("%d", i*10),
			`a.d-block`,
			ctx,
			cancel,
			`href`,
			&sourceNodes,
		)
	}

	// var sourceSitemap []string
	// sitemap_instance := NewSitemap()
	// sitemap_instance.Execute(os.Getenv(`SITEMAP`), &sourceSitemap)

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

	ShuffleSlice(sourceNodes[`NodesValue`].([]string))
	chunked := ChunkSlice(sourceNodes[`NodesValue`], 4)

	for _, i := range chunked.([][]string) {

		var sourceElements []interface{}
		selectors := make(map[string]string)
		selectors[`title`] = `h1.andesbold`
		selectors[`description`] = `p.color-default`
		selectors[`updated`] = `div.metadata-time`
		for _, v := range i {
			threads <- struct{}{}
			elements_instance := NewElements()
			ctx, cancel := InitDriver()
			go elements_instance.Execute(
				`https://datacatalog.worldbank.org`+v,
				selectors,
				ctx,
				cancel,
				&sourceElements,
			)
		}

		for _, i := range sourceElements {

			tmp, ok := i.(map[string]string)
			if !ok {
				continue
			}

			detail := Detail{}
			detail.URL = tmp[`URL`]
			detail.Title = tmp[`title`]
			detail.Describe = tmp[`description`]
			info := Info{}
			info.Updated = tmp[`updated`]
			detail.Info = info

			// if detail.URL == "" && detail.Title == "" {
			// 	continue
			// }

			json_res, _ := json.Marshal(&detail)

			WriteFile(`results`, json_res)
		}

	}

	for {
		if len(threads) == 0 {
			break
		}
	}

	upload_instance := NewUpload()
	upload_instance.Execute(`results`, os.Getenv(`INFO`))

}
