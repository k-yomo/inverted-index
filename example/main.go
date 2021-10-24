package main

import (
	"math/rand"
	"time"

	"github.com/k-yomo/inverted-index/index"
	"github.com/k-yomo/inverted-index/schema"
	"github.com/k0kubun/pp/v3"
)

func main() {
	docs := []index.Document{
		{ID: 1, Text: "there is a white cat"},
		{ID: 2, Text: "black hair cat"},
		{ID: 3, Text: "black cat"},
		{ID: 4, Text: "white dog"},
		{ID: 5, Text: "blue cat"},
		{ID: 6, Text: "black tiger"},
		{ID: 7, Text: "white hair dog"},
	}

	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(docs), func(i, j int) { docs[i], docs[j] = docs[j], docs[i] })

	schemaBuilder := schema.NewBuilder()
	schemaBuilder.AddTextField("title")
	schemaBuilder.AddTextField("description")

	idx, err := index.NewBuilder(schemaBuilder.Build()).CreateInDir("tmp")
	if err != nil {
		panic(err)
	}

	pp.Println(idx)

	//wg := sync.WaitGroup{}
	//for _, doc := range docs {
	//	doc := doc
	//	wg.Add(1)
	//	go func() {
	//		defer wg.Done()
	//		idx.AddDoc(doc)
	//	}()
	//}
	//wg.Wait()
	//idx.DeleteDoc(3)
	//
	//_ = idx.Search("black cat")
	//hits := idx.Search("black cat")
	//pp.Println(hits)
}
