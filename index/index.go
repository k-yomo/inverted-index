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
	invertedIndex map[string]*docPostings
	totalDocNum   int
	totalTokenNum int

	mu *sync.Mutex
}

type docInfo struct {
	tokenNum     int
	uniqueTokens []string
}

func NewIndex(analyzer analyzer.Analyzer) *Index {
	return &Index{
		analyzer:      analyzer,
		forwardIndex:  make(map[int]*docInfo),
		invertedIndex: make(map[string]*docPostings),
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
			idx.invertedIndex[token] = &docPostings{}
		}
		docPos := &docPosting{
			docID:         doc.ID,
			termFrequency: float64(len(postingList)) / float64(len(tokens)),
			postings:      postingList,
		}
		idx.invertedIndex[token].InsertDocPosting(docPos)
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
		dp, ok := idx.invertedIndex[token]
		if !ok {
			continue
		}
		dp.DeleteDocPosting(docID)
	}

	idx.totalTokenNum -= di.tokenNum
	idx.totalDocNum--
	delete(idx.forwardIndex, docID)
}

type phrasePosition struct {
	docID int
	from  int
	to    int
}

type tokenPosition struct {
	docID    int
	position int
}

type Hit struct {
	DocID int
	Score float64
}

func (idx *Index) Search(phrase string) []*Hit {
	tokens := idx.analyzer.Analyze(phrase)
	docScoreMap := make(map[int]float64)
	for _, token := range tokens {
		if dps, ok := idx.invertedIndex[token]; ok {
			for _, dp := range *dps {
				docScoreMap[dp.docID] += idx.okapiBM25(token, dp.docID, dp.termFrequency)
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
	pos := &tokenPosition{docID: 0, position: 0}
	for {
		phrasePos := idx.nextPhrase(tokens, pos)
		if phrasePos.docID == Inf {
			break
		}
		if len(docIDs) == 0 || docIDs[len(docIDs)-1] != phrasePos.docID {
			docIDs = append(docIDs, phrasePos.docID)
		}
		pos = &tokenPosition{
			docID:    phrasePos.docID,
			position: phrasePos.to,
		}
	}
	return docIDs
}

func (idx *Index) nextPhrase(tokens []string, pos *tokenPosition) *phrasePosition {
	tokenNum := len(tokens)
	if tokenNum == 0 {
		return &phrasePosition{
			docID: Inf,
			from:  Inf,
			to:    Inf,
		}
	}

	from := idx.next(tokens[0], pos)
	if from.docID == Inf {
		return &phrasePosition{
			docID: Inf,
			from:  Inf,
			to:    Inf,
		}
	}
	to := &tokenPosition{
		docID:    from.docID,
		position: from.position,
	}
	for i := 1; i < len(tokens); i++ {
		to = idx.next(tokens[i], to)
	}
	if to.docID == from.docID && to.position-from.position == tokenNum-1 {
		return &phrasePosition{
			docID: from.docID,
			from:  from.position,
			to:    to.position,
		}
	}

	return idx.nextPhrase(tokens, from)
}

func (idx *Index) next(token string, pos *tokenPosition) *tokenPosition {
	infPos := &tokenPosition{docID: Inf, position: Inf}

	dps, ok := idx.invertedIndex[token]
	if !ok {
		return infPos
	}
	docPostingList := *dps
	if firstDoc := docPostingList[0]; pos.docID < docPostingList[0].docID {
		return &tokenPosition{docID: firstDoc.docID, position: firstDoc.postings[0]}
	}

	last := docPostingList[len(docPostingList)-1]
	if pos.docID > last.docID || pos.docID == last.docID && pos.position >= last.Last() {
		return infPos
	}

	docPos := dps.NextDocIndex(pos.docID - 1)
	for docPos < len(docPostingList) {
		tokenPos := docPostingList[docPos].NextIndex(pos.position)
		if tokenPos == Inf {
			docPos += 1
			pos.position = NegativeInf
			continue
		}

		return &tokenPosition{
			docID:    docPostingList[docPos].docID,
			position: tokenPos,
		}
	}
	return infPos
}

func (idx *Index) averageDocumentLength() float64 {
	return float64(idx.totalTokenNum) / float64(idx.totalDocNum)
}

// idf calculates TF-IDF's IDF value
func (idx *Index) idf(token string) float64 {
	dps, ok := idx.invertedIndex[token]
	df := 1.0
	if ok {
		df += float64(len(*dps))
	}
	return 1.0 + math.Log2(float64(idx.totalDocNum)/df)
}

// okapiBM25IDF calculates Okapi BM25's IDF value
func (idx *Index) okapiBM25IDF(token string) float64 {
	dps, ok := idx.invertedIndex[token]
	var df float64
	if ok {
		df += float64(len(*dps))
	}
	return 1.0 + math.Log2(1+(float64(idx.totalDocNum)-df+0.5)/(df+0.5))
}

func (idx *Index) okapiBM25(token string, docID int, tf float64) float64 {
	k := 2.0
	b := 0.75
	idf := idx.okapiBM25IDF(token)
	dl := float64(idx.forwardIndex[docID].tokenNum)
	avgDL := idx.averageDocumentLength()
	return idf * ((tf * (k + 1)) / (tf + k*(1-b+b*(dl/avgDL))))
}
