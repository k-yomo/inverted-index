package fieldnorm

import "github.com/k-yomo/inverted-index/schema"

type FieldNormsWriter struct {
	fields           []schema.Field
	fieldnormsBuffer [][]uint8
}

func NewFieldNormsWriter(s *schema.Schema) *FieldNormsWriter {
	return &FieldNormsWriter{
		fields:           s.Fields,
		fieldnormsBuffer: nil,
	}
}
