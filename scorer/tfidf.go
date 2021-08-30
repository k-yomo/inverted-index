package scorer

import "math"

type TFIDFScorer struct {
	idf float64
}

func NewTFIDFScorer(totalDocNum int, documentFrequency int) *TFIDFScorer {
	df := 1.0 + float64(documentFrequency)
	return &TFIDFScorer{
		idf: 1.0 + math.Log2(float64(totalDocNum)/df),
	}
}

func (t *TFIDFScorer) Score(termFrequency float64) float64 {
	return termFrequency * t.idf
}
