package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"

	sitemap "github.com/oxffaa/gopher-parse-sitemap"
)

type Sitemap struct {
	SitemapURL []string
}

/*
	params[0]url, params[1]source(return), params[2]regexp(optional)
*/
func (s *Sitemap) Execute(params ...interface{}) (interface{}, error) {

	defer s.Getter(params[1])
	if len(params) == 2 {
		sitemapurl := s.readSitemap(fmt.Sprintf("%v", params[0]))
		s.SitemapURL = make([]string, len(sitemapurl))
		copy(s.SitemapURL, sitemapurl)
		return sitemapurl, nil
	} else if len(params) == 3 {
		sitemapurl := s.readSitemap(fmt.Sprintf("%v", params[0]))
		parsedurl := s.parseSitemap(fmt.Sprintf("%v", params[1]), sitemapurl)
		s.SitemapURL = make([]string, len(parsedurl))
		copy(s.SitemapURL, parsedurl)
		return parsedurl, nil
	}
	return nil, fmt.Errorf("Not enough params")
}

func (s *Sitemap) readSitemap(url string) []string {

	resp, _ := http.Get(url)
	defer resp.Body.Close()

	out, _ := os.Create("sitemap")
	defer out.Close()

	io.Copy(out, resp.Body)

	file, _ := os.Open("sitemap")
	defer file.Close()

	var sitemapurl []string

	sitemap.ParseFromFile("./sitemap", func(e sitemap.Entry) error {
		sitemapurl = append(sitemapurl, e.GetLocation())
		return nil
	})

	return sitemapurl
}

func (s *Sitemap) parseSitemap(reg string, sitemapurl []string) []string {

	var parsedurl []string
	for _, v := range sitemapurl {
		reg := regexp.MustCompile(reg)
		url := reg.FindStringSubmatch(v)

		if len(url) != 0 {
			parsedurl = append(parsedurl, url[1])
		}
	}

	return parsedurl
}

func NewSitemap() Icommand {
	return &Sitemap{}
}

/*
	[]string
*/
func (s *Sitemap) Getter(source interface{}) {

	sourceConver, ok := source.(*[]string)
	if !ok {
		return
	}

	*sourceConver = append(*sourceConver, v...)
}
