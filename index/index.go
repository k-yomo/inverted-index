package index

import (
	"github.com/k-yomo/inverted-index/analyzer"
	"math"
	"sort"
	"sync"
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
	analyzer      analyzer.Analyzer
	forwardIndex  map[int]*docInfo
	invertedIndex map[string]*postingInfo
	totalDocNum   int
	totalTokenNum int

	mu *sync.Mutex
}

type docInfo struct {
	tokenNum     int
	uniqueTokens []string
}

type postingInfo struct {
	postings *postings
}

func NewIndex(analyzer analyzer.Analyzer) *Index {
	return &Index{
		analyzer:      analyzer,
		forwardIndex:  make(map[int]*docInfo),
		invertedIndex: make(map[string]*postingInfo),
		totalDocNum:   0,
		totalTokenNum: 0,
		mu:            &sync.Mutex{},
	}
}

func (idx *Index) AddDocs(documents ...Document) {
	for _, doc := range documents {
		idx.AddDoc(doc)
	}
}

func (idx *Index) AddDoc(doc Document) {
	idx.mu.Lock()
	defer idx.mu.Unlock()

	if _, ok := idx.forwardIndex[doc.ID]; ok {
		idx.deleteDoc(doc.ID)
	}

	postingMap := make(map[string][]int)
	tokens := idx.analyzer.Analyze(doc.Text)
	for i, token := range tokens {
		postingMap[token] = append(postingMap[token], i)
	}

	var uniqueTokens []string
	for token, postingList := range postingMap {
		uniqueTokens = append(uniqueTokens, token)
		if idx.invertedIndex[token] == nil {
			idx.invertedIndex[token] = &postingInfo{
				postings: &postings{},
			}
		}
		pos := &posting{
			docID:         doc.ID,
			termFrequency: float64(len(postingList)) / float64(len(tokens)),
			postings:      postingList,
		}
		idx.invertedIndex[token].postings.InsertPosting(pos)
	}

	idx.forwardIndex[doc.ID] = &docInfo{
		tokenNum:     len(tokens),
		uniqueTokens: uniqueTokens,
	}

	idx.totalDocNum++
	idx.totalTokenNum += len(tokens)
}

func (idx *Index) DeleteDoc(docID int) {
	idx.mu.Lock()
	defer idx.mu.Unlock()

	idx.deleteDoc(docID)
}

// deleteDoc deletes doc
// lock must be taken to use this function since this is not thread-safe
func (idx *Index) deleteDoc(docID int) {
	di, ok := idx.forwardIndex[docID]
	if !ok {
		return
	}

	for _, token := range di.uniqueTokens {
		pi, ok := idx.invertedIndex[token]
		if !ok {
			continue
		}
		pi.postings.DeletePosting(docID)
	}

	idx.totalTokenNum -= di.tokenNum
	idx.totalDocNum--
	delete(idx.forwardIndex, docID)
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

type Hit struct {
	DocID int
	Score float64
}

func (idx *Index) Search(phrase string) []*Hit {
	tokens := idx.analyzer.Analyze(phrase)
	docScoreMap := make(map[int]float64)
	for _, token := range tokens {
		if postingInfo, ok := idx.invertedIndex[token]; ok {
			for _, posting := range *postingInfo.postings {
				docScoreMap[posting.docID] += idx.okapiBM25(token, posting.docID, posting.termFrequency)
			}
		}
	}

	hits := make([]*Hit, 0, len(docScoreMap))
	for docID, score := range docScoreMap {
		hits = append(hits, &Hit{DocID: docID, Score: score})
	}

	sort.Slice(hits, func(i, j int) bool { return hits[i].Score > hits[j].Score })

	return hits
}

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
	pi, ok := idx.invertedIndex[token]
	if !ok {
		return InfPosition
	}
	postingList := *pi.postings
	if firstDoc := postingList[0]; position.docID < postingList[0].docID {
		return &Position{docID: firstDoc.docID, tokenPosition: firstDoc.postings[0]}
	}

	lastDoc := postingList[len(postingList)-1]
	if position.docID > lastDoc.docID || position.docID == lastDoc.docID && position.tokenPosition >= lastDoc.Last() {
		return InfPosition
	}

	docPos := pi.postings.NextDocIndex(position.docID - 1)
	for docPos < len(postingList) {
		tokenPos := postingList[docPos].NextIndex(position.tokenPosition)
		if tokenPos == Inf {
			docPos += 1
			position.tokenPosition = NegativeInf
			continue
		}

		return &Position{
			docID:         postingList[docPos].docID,
			tokenPosition: tokenPos,
		}
	}
	return &Position{
		docID:         Inf,
		tokenPosition: Inf,
	}
}

func (idx *Index) averageDL() float64 {
	return float64(idx.totalTokenNum) / float64(idx.totalDocNum)
}

// idf calculates TF-IDF's IDF value
func (idx *Index) idf(token string) float64 {
	postingInfo, ok := idx.invertedIndex[token]
	df := 1.0
	if ok {
		df += float64(len(*postingInfo.postings))
	}
	return 1.0 + math.Log2(float64(idx.totalDocNum)/df)
}

// okapiBM25IDF calculates Okapi BM25's IDF value
func (idx *Index) okapiBM25IDF(token string) float64 {
	postingInfo, ok := idx.invertedIndex[token]
	var df float64
	if ok {
		df += float64(len(*postingInfo.postings))
	}
	return 1.0 + math.Log2(1+(float64(idx.totalDocNum)-df+0.5)/(df+0.5))
}

func (idx *Index) okapiBM25(token string, docID int, tf float64) float64 {
	k := 2.0
	b := 0.75
	idf := idx.okapiBM25IDF(token)
	dl := float64(idx.forwardIndex[docID].tokenNum)
	return idf * ((tf * (k + 1)) / (tf + k*(1-b+b*(dl/idx.averageDL()))))
}
