package index

import (
	"regexp"
	"strings"
)

const (
	Inf         = int(^uint(0) >> 1) // use max int as Inf
	NegativeInf = -Inf - 1           // use min int as -Inf
)

type Index struct {
	index map[string][]int
}

func NewIndex(document string) *Index {
	terms := strings.Split(document, " ")
	index := make(map[string][]int)
	for i, term := range terms {
		term = regexp.MustCompile("[\\W]+").ReplaceAllString(strings.ToLower(term), "")
		index[term] = append(index[term], i)
	}
	return &Index{
		index: index,
	}
}

type Range struct {
	From int
	To   int
}

func (idx *Index) NextPhrase(phrase string, position int) *Range {
	terms := strings.Split(phrase, " ")
	v := position
	for i := 0; i < len(terms); i++ {
		v = idx.Next(terms[i], v)
	}
	if v == Inf {
		return &Range{From: Inf, To: Inf}
	}
	u := v
	for i := len(terms) - 2; i >= 0; i-- {
		u = idx.Prev(terms[i], u)
	}
	if v-u == len(terms)-1 {
		return &Range{From: u, To: v}
	}
	return idx.NextPhrase(phrase, u)
}

func (idx *Index) First(term string) int {
	postings := idx.index[term]
	if len(postings) == 0 {
		return NegativeInf
	}
	return postings[0]
}

func (idx *Index) Last(term string) int {
	postings := idx.index[term]
	if len(postings) == 0 {
		return Inf
	}
	return postings[0]
}

func (idx *Index) Prev(term string, current int) int {
	postings := idx.index[term]
	// TODO: refactor with binary search
	for i := len(postings) - 1; i >= 0; i-- {
		if postings[i] < current {
			return postings[i]
		}
	}
	return NegativeInf
}

func (idx *Index) Next(term string, current int) int {
	postings := idx.index[term]
	// TODO: refactor with binary search
	for i, pos := range postings {
		if pos > current {
			return postings[i]
		}
	}
	return Inf
}
