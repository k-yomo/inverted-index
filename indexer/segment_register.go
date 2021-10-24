package indexer

import "github.com/k-yomo/inverted-index/index"

type SegmentRegister struct {
	segmentStatus map[index.SegmentID]*SegmentEntry
}

func newSegmentRegister() *SegmentRegister {
	return &SegmentRegister{
		segmentStatus: make(map[index.SegmentID]*SegmentEntry),
	}
}

func newSegmentRegisterFromSegmentMetas(segmentMetas []*index.SegmentMeta) *SegmentRegister {
	segmentStatus := make(map[index.SegmentID]*SegmentEntry)
	for _, segmentMeta := range segmentMetas {
		segmentStatus[segmentMeta.SegmentID] = NewSegmentEntry(segmentMeta)
	}

	return &SegmentRegister{
		segmentStatus: segmentStatus,
	}
}
