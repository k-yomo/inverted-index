package index

import (
	"regexp"
	"strings"
)

const (
	Inf         = int(^uint(0) >> 1) // use max int as ∞
	NegativeInf = -Inf - 1           // use min int as -∞
)

type Index struct {
	PostingMap map[string][]int
}

func NewIndex(document string) *Index {
	terms := strings.Split(document, " ")
	index := make(map[string][]int)
	for i, term := range terms {
		term = regexp.MustCompile("[\\W]+").ReplaceAllString(strings.ToLower(term), "")
		if term != "" {
			index[term] = append(index[term], i)
		}
	}
	return &Index{
		PostingMap: index,
	}
}

type Range struct {
	From int
	To   int
}

func (idx *Index) NextPhrase(phrase string, position int) *Range {
	terms := strings.Split(phrase, " ")
	termNum := len(terms)

	v := position
	for i := 0; i < len(terms); i++ {
		v = idx.Next(terms[i], v)
	}
	if v == Inf {
		return &Range{From: Inf, To: Inf}
	}

	u := v
	for i := termNum-2; i >= 0; i-- {
		u = idx.Prev(terms[i], u)
	}

	if v-u == termNum-1 {
		return &Range{From: u, To: v}
	}
	return idx.NextPhrase(phrase, u)
}

func (idx *Index) First(term string) int {
	postings := idx.PostingMap[term]
	if len(postings) == 0 {
		return NegativeInf
	}
	return postings[0]
}

func (idx *Index) Last(term string) int {
	postings := idx.PostingMap[term]
	if len(postings) == 0 {
		return Inf
	}
	return postings[len(postings)-1]
}

func (idx *Index) Prev(term string, current int) int {
	postings, exist := idx.PostingMap[term]
	if !exist || current <= postings[0] {
		return NegativeInf
	}
	if current > postings[len(postings)-1] {
		return idx.Last(term)
	}

	ok, ng := 0, len(postings)-1
	for ng-ok > 1 {
		mid := (ok + ng) / 2
		if postings[mid] < current {
			ok = mid
		} else {
			ng = mid
		}
	}

	return postings[ok]
}

func (idx *Index) Next(term string, current int) int {
	postings, exist := idx.PostingMap[term]
	if !exist || current >= postings[len(postings)-1] {
		return Inf
	}
	if current < postings[0] {
		return idx.First(term)
	}

	ok, ng := len(postings)-1, 0
	for ok-ng > 1 {
		mid := (ok + ng) / 2
		if postings[mid] > current {
			ok = mid
		} else {
			ng = mid
		}
	}
	return postings[ok]
}
