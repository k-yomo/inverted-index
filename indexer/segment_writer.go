package indexer

import (
	"github.com/k-yomo/inverted-index/analyzer"
	"github.com/k-yomo/inverted-index/fieldnorm"
	"github.com/k-yomo/inverted-index/index"
	"github.com/k-yomo/inverted-index/internal/opstamp"
	"github.com/k-yomo/inverted-index/postings"
	"github.com/k-yomo/inverted-index/schema"
)

type SegmentWriter struct {
	maxDoc             uint32
	multifieldPostings *postings.MultiFieldPostingsWriter
	fieldnormsWriter   *fieldnorm.FieldNormsWriter
	docOpstamps        []opstamp.Opstamp
	analyzer           analyzer.Analyzer
}

func newSegmentWriter(memoryBudget int, segment *index.Segment, schema *schema.Schema) *SegmentWriter {
	return &SegmentWriter{
		maxDoc:             0,
		multifieldPostings: postings.NewMultiFieldPostingsWriter(schema),
		docOpstamps:        nil,
		analyzer:           segment.Index.Analyzer,
	}
}

func (s *SegmentWriter) addDocument(addOperation *AddOperation, sc *schema.Schema) error {
	docID := s.maxDoc
	doc := addOperation.document
	s.docOpstamps = append(s.docOpstamps, addOperation.opstamp)
	for _, fieldAndFieldValues := range doc.SortedFieldValues() {
		fieldEntry := sc.Fields[fieldAndFieldValues.Field]
		switch fieldEntry.FieldType {
		case schema.FieldTypeText:
			var tokensList [][]string
			var offsets []int
			var totalOffset int

			for _, fieldValue := range fieldAndFieldValues.FieldValues {
				switch v := fieldValue.Value.(type) {
				case string:
					offsets = append(offsets, totalOffset)
					totalOffset += len(v)
					tokensList = append(tokensList, s.analyzer.Analyze(v))
				}
			}

			var tokenNum int
			if len(tokensList) > 0 {
				s.multifieldPostings.IndexText(docID, fieldAndFieldValues.Field, tokensList)
			}

		}
	}
}
