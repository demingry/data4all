package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"
)

var (
	threads chan struct{}
	t       int64
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

	// params := strings.Split(os.Getenv("REQUEST_PARAMS"), "-")
	// start, _ := strconv.Atoi(params[0])
	// end, _ := strconv.Atoi(params[1])

	// sourceNodes := make(map[string]interface{})
	// for i := start; i < end; i++ {
	// 	fmt.Println("Scraped Page: " + fmt.Sprintf("%d", i))
	// 	threads <- struct{}{}
	// 	ctx, cancel := InitDriver()
	// 	nodes_instance := NewNodes()
	// 	go nodes_instance.Execute(
	// 		`https://data.gov.uk/search?filters%5Btopic%5D=Society&page=`+fmt.Sprintf("%d", i),
	// 		`.govuk-heading-m a`,
	// 		ctx,
	// 		cancel,
	// 		`href`,
	// 		&sourceNodes,
	// 	)
	// }

	var sourceSitemap []string
	sitemap_instance := NewSitemap()
	sitemap_instance.Execute(os.Getenv(`SITEMAP`), &sourceSitemap)

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

	chunked := ChunkSlice(sourceSitemap, 4)
	fmt.Println(len(chunked.([][]string)))

	for _, i := range chunked.([][]string) {
		var sourcePage []string
		for _, v := range i {
			threads <- struct{}{}
			fmt.Println(v)
			page_instance := NewPage(`PageFromDriver`)
			ctx, cancel := InitDriver()
			go page_instance.Execute(
				v,
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

		var sourceSoup []string
		for _, v := range sourcePage {
			threads <- struct{}{}
			soup_instance := NewSoup()
			go soup_instance.Execute(
				v,
				[]string{`script`, `type`, `application/ld+json`},
				&sourceSoup,
			)
		}

		for {
			if len(threads) == 0 {
				break
			}
		}

		for _, i := range sourceSoup {

			auto := &AutoGenerated{}
			json.Unmarshal([]byte(i), &auto)

			detail := Detail{}
			detail.URL = auto.URL
			detail.Title = auto.Name
			detail.Describe = auto.Description
			info := Info{}
			if len(auto.Author) != 0 {
				info.Publisher = auto.Author[0].Name
			}
			info.Created = auto.DatePublished
			info.Identifier = auto.Identifier
			detail.Info = info

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
