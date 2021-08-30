package index

type postingList struct {
	docPostings       *docPostings
	documentFrequency int
}

type docPostings []*docPosting

type docPosting struct {
	docID         int
	termFrequency float64
	positions     []int
	positionCache int
}

func (ps docPostings) PrevDocIndex(docID int) int {
	if len(ps) == 0 || docID <= ps[0].docID {
		return NegativeInf
	}
	if docID == Inf {
		return len(ps) - 1
	}

	ok, ng := 0, len(ps)-1
	for ng-ok > 1 {
		mid := (ok + ng) / 2
		if ps[mid].docID < docID {
			ok = mid
		} else {
			ng = mid
		}
	}

	return ok
}

func (ps docPostings) NextDocIndex(docID int) int {
	if len(ps) == 0 || docID >= ps[len(ps)-1].docID {
		return Inf
	}
	if docID < ps[0].docID {
		return 0
	}

	high, low := len(ps)-1, 0
	for high-low > 1 {
		mid := (low + high) / 2
		if ps[mid].docID > docID {
			high = mid
		} else {
			low = mid
		}
	}
	return high
}

func (p *postingList) InsertDocPosting(pos *docPosting) {
	list := *p.docPostings
	var updated docPostings

	nextIdx := list.NextDocIndex(pos.docID)
	if nextIdx == Inf {
		updated = append(list, pos)
	} else {
		updated = append(
			list[:nextIdx],
			append([]*docPosting{pos}, list[nextIdx:]...)...,
		)
	}
	p.docPostings = &updated
	p.documentFrequency++
}

func (p *postingList) DeleteDocPosting(docID int) {
	list := *p.docPostings
	docPos := list.NextDocIndex(docID - 1)
	if docPos >= len(list) || list[docPos].docID != docID {
		return
	}
	*p.docPostings = append(list[:docPos], list[docPos+1:]...)
	p.documentFrequency--
}

func (p *docPosting) First() int {
	if len(p.positions) == 0 {
		return NegativeInf
	}
	return p.positions[0]
}

func (p *docPosting) Last() int {
	if len(p.positions) == 0 {
		return Inf
	}
	return p.positions[len(p.positions)-1]
}

func (p *docPosting) PrevIndex(current int) int {
	postings := p.positions
	if current <= postings[0] {
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

func (p *docPosting) NextIndex(current int) int {
	last := p.Last()
	postings := p.positions
	if current >= last {
		return Inf
	}
	if current < postings[0] {
		p.positionCache = 0
		return postings[0]
	}

	low := 0
	if cache := p.positionCache; cache > 0 && postings[cache-1] <= current {
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
	p.positionCache = nextIndex
	return postings[nextIndex]
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
