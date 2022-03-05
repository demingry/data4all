package main

import "time"

type Icommand interface {
	Execute(...interface{}) (interface{}, error)
}

type IGetter interface {
	Getter(interface{})
}

//Result Structure

type Detail struct {
	URL      string `json:"URL"`
	Title    string `json:"Title"`
	Describe string `json:"Describe"`
	Info     Info   `json:"Info"`
}

type Info struct {
	Publisher  string `json:"Publisher"`
	Created    string `json:"Created"`
	Updated    string `json:"Updated"`
	Identifier string `json:Identifier`
}

//Json Unmarshal Structure(Autogenerated Customized)

type AutoGenerated struct {
	Context          string      `json:"@context"`
	Type             string      `json:"@type"`
	Name             string      `json:"name"`
	Description      string      `json:"description"`
	URL              string      `json:"url"`
	SameAs           string      `json:"sameAs"`
	Version          string      `json:"version"`
	Keywords         string      `json:"keywords"`
	Publisher        interface{} `json:"publisher"`
	DatePublished    time.Time   `json:"datePublished"`
	MainEntityOfPage struct {
		Type string `json:"@type"`
		ID   string `json:"@id"`
	} `json:"mainEntityOfPage"`
	IncludedInDataCatalog string `json:"includedInDataCatalog"`
	License               struct {
		Type string `json:"@type"`
		URL  string `json:"url"`
		Text string `json:"text"`
	} `json:"license"`
	Author struct {
		Type string `json:"@type"`
		Name string `json:"name"`
	} `json:"author"`
	Citation struct {
		Type          string    `json:"@type"`
		Text          string    `json:"text"`
		Headline      string    `json:"headline"`
		DatePublished time.Time `json:"datePublished"`
		URL           string    `json:"url"`
	} `json:"citation"`
}
