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
			`https://data.gov.au/search?page=`+fmt.Sprintf("%d", i),
			`.dataset-summary-title a`,
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

	fmt.println("this")

	chunked := ChunkSlice(sourceNodes["NodesValue"], 4)

	for _, v := range chunked.([][]string) {

		var sourceElements []interface{}
		selectors := make(map[string]string)
		selectors[`URL`] = `.au-breadcrumbs li:nth-of-type(3) a`
		selectors[`title`] = `h1`
		selectors[`description`] = `.no-print p`
		selectors[`created`] = `span[itemprop='dateCreated']`
		selectors[`updated`] = `span[itemprop='dateModified']`
		selectors[`publisher`] = `[itemprop='url'] span`
		for _, i := range v {
			threads <- struct{}{}
			fmt.Println(v)
			elements_instance := NewElements()
			ctx, cancel := InitDriver()
			go elements_instance.Execute(
				`https://data.gov.au`+i,
				selectors,
				ctx,
				cancel,
				&sourceElements,
			)
		}

		for {
			if len(threads) == 0 {
				break
			}
		}

		for _, i := range sourceElements {

			tmp, ok := i.(map[string]string)
			if !ok {
				continue
			}

			detail := Detail{}
			detail.URL = tmp["URL"]
			detail.Title = tmp["title"]
			detail.Describe = tmp["describe"]
			info := Info{}
			info.Publisher = tmp["publisher"]
			info.Created = tmp["created"]
			info.Updated = tmp["updated"]
			detail.Info = info

			if detail.URL == "" && detail.Title == "" {
				continue
			}

			json_res, _ := json.Marshal(&detail)

			WriteFile(`results`, json_res)
		}

		// var sourceSoup []string
		// for _, v := range sourcePage {
		// 	threads <- struct{}{}
		// 	soup_instance := NewSoup()
		// 	go soup_instance.Execute(
		// 		v,
		// 		[]string{`script`, `type`, `application/ld+json`},
		// 		&sourceSoup,
		// 	)
		// }

	}

	for {
		if len(threads) == 0 {
			break
		}
	}

	upload_instance := NewUpload()
	upload_instance.Execute(`results`, os.Getenv(`INFO`))

}
