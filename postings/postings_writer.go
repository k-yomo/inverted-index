package postings

import "github.com/k-yomo/inverted-index/schema"

type MultiFieldPostingsWriter struct {
	schema *schema.Schema
}

func NewMultiFieldPostingsWriter(s *schema.Schema) *MultiFieldPostingsWriter {
	return &MultiFieldPostingsWriter{
		schema: s,
	}
}

func (m *MultiFieldPostingsWriter) IndexText(docID uint32, field schema.Field, tokensList [][]string) {

}
