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
	PostingMap   map[string][]int
	postingCache map[string]int
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
		postingCache: make(map[string]int),
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
	for i := termNum - 2; i >= 0; i-- {
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
		return postings[len(postings)-1]
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
	last := idx.Last(term)
	postings, exist := idx.PostingMap[term]
	if !exist || current >= last {
		return Inf
	}
	if current < postings[0] {
		idx.setPostingCache(term, 0)
		return postings[0]
	}

	low := 0
	if cache := idx.postingCache[term]; cache > 0 && postings[cache-1] <= current {
		low = cache - 1
	}

	jump := 1
	high := low + jump
	for high < len(postings)-1 && postings[high] <= current {
		low = high
		jump = 2 * jump
		high = low + jump
	}
	if high > last {
		high = last
	}

	nextIndex := binarySearch(postings, low, high, current)
	idx.setPostingCache(term, nextIndex)
	return postings[nextIndex]
}

func (idx *Index) setPostingCache(term string, i int) {
	idx.postingCache[term] = i
}

func binarySearch(postings []int, low, high, current int) int {
	for high-low > 1 {
		mid := (low + high) / 2
		if postings[mid] > current {
			high = mid
		} else {
			low = mid
		}
	}
	return high
}
