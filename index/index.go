package index

import (
	"github.com/k-yomo/inverted-index/analyzer"
)

const (
	Inf         = int(^uint(0) >> 1) // use max int as ∞
	NegativeInf = -Inf - 1           // use min int as -∞
)

type Document struct {
	ID   int
	Text string
}

type Index struct {
	analyzer            analyzer.Analyzer
	tokenPostingInfoMap map[string]*postingInfo
	docNum              int
}

type postingInfo struct {
	postings     postings
	docFrequency float64
}

func NewIndex(analyzer analyzer.Analyzer, documents ...Document) *Index {
	tokenPostingInfoMap := make(map[string]*postingInfo)
	for _, doc := range documents {
		postingMap := make(map[string][]int)
		tokens := analyzer.Analyze(doc.Text)
		for i, term := range tokens {
			postingMap[term] = append(postingMap[term], i)
		}
		for term, postings := range postingMap {
			if tokenPostingInfoMap[term] == nil {
				tokenPostingInfoMap[term] = &postingInfo{}
			}
			tokenPostingInfoMap[term].postings = append(tokenPostingInfoMap[term].postings, &posting{
				docID:         doc.ID,
				termFrequency: float64(len(postings)) / float64(len(tokens)),
				postings:      postings,
			})
		}
	}

	for _, postingInfo := range tokenPostingInfoMap {
		postingInfo.docFrequency = float64(len(postingInfo.postings)) / float64(len(documents))
	}

	return &Index{
		analyzer:            analyzer,
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
	docID         int
	tokenPosition int
}

var InfPosition = &Position{docID: Inf, tokenPosition: Inf}

func (idx *Index) PhraseSearch(phrase string) []int {
	var docIDs []int

	tokens := idx.analyzer.Analyze(phrase)
	position := &Position{docID: 0, tokenPosition: 0}
	for {
		pos := idx.nextPhrase(tokens, position)
		if pos.docID == Inf {
			break
		}
		if len(docIDs) == 0 || docIDs[len(docIDs)-1] != pos.docID {
			docIDs = append(docIDs, pos.docID)
		}
		position = &Position{
			docID:         pos.docID,
			tokenPosition: pos.positionRange.To,
		}
	}
	return docIDs
}

func (idx *Index) nextPhrase(tokens []string, position *Position) *PhrasePosition {
	tokenNum := len(tokens)
	if tokenNum == 0 {
		return &PhrasePosition{
			docID:         Inf,
			positionRange: &Range{From: Inf, To: Inf},
		}
	}

	from := idx.next(tokens[0], position)
	if from.docID == Inf {
		return &PhrasePosition{
			docID:         Inf,
			positionRange: &Range{From: Inf, To: Inf},
		}
	}
	to := from
	for i := 1; i < len(tokens); i++ {
		to = idx.next(tokens[i], to)
	}
	if to.docID == from.docID && to.tokenPosition-from.tokenPosition == tokenNum-1 {
		return &PhrasePosition{
			docID:         from.docID,
			positionRange: &Range{From: from.tokenPosition, To: to.tokenPosition},
		}
	}

	return idx.nextPhrase(tokens, from)
}

func (idx *Index) next(token string, position *Position) *Position {
	pi, ok := idx.tokenPostingInfoMap[token]
	if !ok {
		return InfPosition
	}
	if firstDoc := pi.postings[0]; position.docID < pi.postings[0].docID {
		return &Position{docID: firstDoc.docID, tokenPosition: firstDoc.postings[0]}
	}

	lastDoc := pi.postings[len(pi.postings)-1]
	if position.docID > lastDoc.docID || position.docID == lastDoc.docID && position.tokenPosition >= lastDoc.Last() {
		return InfPosition
	}

	docPos := pi.postings.DocNextIndex(position.docID - 1)
	for docPos < len(pi.postings) {
		tokenPos := pi.postings[docPos].NextIndex(position.tokenPosition)
		if tokenPos == Inf {
			docPos += 1
			position.tokenPosition = NegativeInf
			continue
		}

		return &Position{
			docID:         pi.postings[docPos].docID,
			tokenPosition: tokenPos,
		}
	}
	return &Position{
		docID:         Inf,
		tokenPosition: Inf,
	}
}
