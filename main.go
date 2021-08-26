package main

import (
	"fmt"
	"github.com/k-yomo/inverted-index/analyzer"
	"github.com/k-yomo/inverted-index/index"
)

func main()  {
	docs := []index.Document{
		{ ID: 1, Text: "there is a black cat"},
		{ ID: 2, Text: "black hair cat"},
		{ ID: 3, Text: "black cat"},
	}
	idx := index.NewIndex(&analyzer.EnglishAnalyzer{}, docs...)
	docIDs := idx.PhraseSearch("hair cat")
	fmt.Println(docIDs)
}
