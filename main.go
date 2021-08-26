package main

import (
	"github.com/k-yomo/inverted-index/analyzer"
	"github.com/k-yomo/inverted-index/index"
	"github.com/k0kubun/pp"
)

func main()  {
	docs := []index.Document{
		{ ID: 1, Text: "there is a white cat"},
		{ ID: 2, Text: "black hair cat"},
		{ ID: 3, Text: "black cat"},
		{ ID: 4, Text: "white dog"},
	}

	idx := index.NewIndex(&analyzer.EnglishAnalyzer{}, docs...)
	idx.DeleteDoc(3)
	hits := idx.Search("black cat")

	pp.Println(hits)
}
