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
	postingMap   map[string][]int
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
		postingMap:   index,
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
	if termNum == 0 {
		return &Range{From: Inf, To: Inf}
	}

	from := idx.Next(terms[0], position)
	to := from
	for i := 1; i < len(terms); i++ {
		to = idx.Next(terms[i], to)
	}
	if to == Inf {
		return &Range{From: Inf, To: Inf}
	}

	if to-from == termNum-1 {
		return &Range{From: from, To: to}
	}
	return idx.NextPhrase(phrase, from)
}

func (idx *Index) First(term string) int {
	postings := idx.postingMap[term]
	if len(postings) == 0 {
		return NegativeInf
	}
	return postings[0]
}

func (idx *Index) Last(term string) int {
	postings := idx.postingMap[term]
	if len(postings) == 0 {
		return Inf
	}
	return postings[len(postings)-1]
}

func (idx *Index) Prev(term string, current int) int {
	postings, exist := idx.postingMap[term]
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
	postings, exist := idx.postingMap[term]
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
