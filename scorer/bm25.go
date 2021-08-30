package scorer

import "math"

type BM25Scorer struct {
	idf                      float64
	averageDocumentFrequency float64
}

func NewBM25Scorer(totalDocNum int, totalTokenNum int, documentFrequency int) *BM25Scorer {
	df := float64(documentFrequency)
	idf := 1.0 + math.Log2(1+(float64(totalDocNum)-df+0.5)/(df+0.5))
	return &BM25Scorer{
		idf:                      idf,
		averageDocumentFrequency: float64(totalTokenNum) / float64(totalDocNum),
	}
}

func (bs *BM25Scorer) Score(documentLength int, termFrequency float64) float64 {
	const k = 2.0
	const b = 0.75
	return bs.idf * ((termFrequency * (k + 1)) / (termFrequency + k*(1-b+b*(float64(documentLength)/bs.averageDocumentFrequency))))
}
