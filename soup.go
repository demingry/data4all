package main

import (
	"fmt"

	"github.com/anaskhan96/soup"
)

type Soup struct {
	soupElement string
}

/*
	params[0]HTMLBody, params[1]selector([]string), params[2]source(return)
*/
func (s *Soup) Execute(params ...interface{}) (interface{}, error) {

	defer Finished()
	defer s.Getter(params[2])
	selector, ok := params[1].([]string)
	if !ok {
		return nil, fmt.Errorf("Unexpectable params")
	}
	s.soupParse(fmt.Sprintf("%v", params[0]), selector)
	return nil, nil
}

/*
	selector: type + attrname + attrvalue
*/
func (s *Soup) soupParse(body string, selector []string) string {

	doc := soup.HTMLParse(body)
	soupresult := doc.Find(selector...)
	if soupresult.Error != nil {
		return ""
	}
	return soupresult.Text()
}

func NewSoup() Icommand {
	return &Soup{}
}

func (s *Soup) Getter(source interface{}) {

	mu.Lock()
	sourceConver, ok := source.(*[]string)
	if !ok {
		return
	}

	*sourceConver = append(*sourceConver, s.soupElement)
	mu.Unlock()
}
