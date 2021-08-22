package index

const (
	Inf         = int(^uint(0) >> 1) // use max int as ∞
	NegativeInf = -Inf - 1           // use min int as -∞
)

type Document struct {
	ID   int
	Text string
}

type Index struct {
	tokenPostingInfoMap map[string]*postingInfo
	docNum              int
}

type postingInfo struct {
	postings     postings
	docFrequency float64
}

func NewIndex(documents ...Document) *Index {
	tokenPostingInfoMap := make(map[string]*postingInfo)
	for _, doc := range documents {
		postingMap := make(map[string][]int)
		terms := tokenize(doc.Text)
		for i, term := range terms {
			postingMap[term] = append(postingMap[term], i)
		}
		for term, postings := range postingMap {
			if tokenPostingInfoMap[term] == nil {
				tokenPostingInfoMap[term] = &postingInfo{}
			}
			tokenPostingInfoMap[term].postings = append(tokenPostingInfoMap[term].postings, &posting{
				docID:         doc.ID,
				termFrequency: float64(len(postings)) / float64(len(terms)),
				postings:      postings,
			})
		}
	}

	for _, postingInfo := range tokenPostingInfoMap {
		postingInfo.docFrequency = float64(len(postingInfo.postings)) / float64(len(documents))
	}

	return &Index{
		tokenPostingInfoMap: tokenPostingInfoMap,
		docNum:              len(documents),
	}
}

type PhrasePosition struct {
	docID         int
	positionRange *Range
}

type TermPosition struct {
	docID    int
	position int
}

type Range struct {
	From int
	To   int
}

type Position struct {
	docID    int
	position int
}

func (idx *Index) NextPhrase(phrase string, position *Position) *PhrasePosition {
	terms := tokenize(phrase)
	termNum := len(terms)
	if termNum == 0 {
		return &PhrasePosition{
			docID:         Inf,
			positionRange: &Range{From: Inf, To: Inf},
		}
	}

	from := idx.Next(terms[0], position)
	if from.docID == Inf {
		return &PhrasePosition{
			docID:         Inf,
			positionRange: &Range{From: Inf, To: Inf},
		}
	}
	to := from
	for i := 1; i < len(terms); i++ {
		to = idx.Next(terms[i], to)
	}

	if to.docID == from.docID && to.position-from.position == termNum-1 {
		return &PhrasePosition{
			docID:         from.docID,
			positionRange: &Range{From: from.position, To: to.position},
		}
	}
	return idx.NextPhrase(phrase, from)
}

func (idx *Index) Next(term string, position *Position) *Position {
	if position.docID == Inf {
		return &Position{docID: Inf, position: Inf}
	}

	pi := idx.tokenPostingInfoMap[term]
	for position.docID != Inf {

	}
	if position.docID < pi.postings[0].docID {
		pos := pi.postings[0].Next(position.position)
		if pos == Inf {

		}
	}
	postingInfo, ok := idx.tokenPostingInfoMap[term]
	if !ok {
		return &Position{docID: Inf, position: Inf}
	}
	postingInfo.postings
}
