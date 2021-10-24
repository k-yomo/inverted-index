package index

import "github.com/k-yomo/inverted-index/schema"

type Segment struct {
	Index *Index
	meta  *SegmentMeta
}

func newSegment(idx *Index, meta *SegmentMeta) *Segment {
	return &Segment{
		Index: idx,
		meta:  meta,
	}
}

func (s *Segment) Schema() *schema.Schema {
	return s.Index.schema
}
