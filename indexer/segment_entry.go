package indexer

import "github.com/k-yomo/inverted-index/index"

type SegmentEntry struct {
	meta *index.SegmentMeta
}

func NewSegmentEntry(segmentMeta *index.SegmentMeta) *SegmentEntry {
	return &SegmentEntry{meta: segmentMeta}
}
