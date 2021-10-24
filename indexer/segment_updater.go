package indexer

import (
	"github.com/k-yomo/inverted-index/index"
	"github.com/k-yomo/inverted-index/internal/opstamp"
)

type SegmentUpdater struct {
	indexMeta      *index.IndexMeta
	index          *index.Index
	segmentManager *SegmentManager
}

func NewSegmentUpdater(idx *index.Index, stamper *opstamp.Stamper) (*SegmentUpdater, error) {
	indexMeta, err := idx.LoadMetas()
	if err != nil {
		return nil, err
	}
	segmentManager := NewSegmentManager(indexMeta.Segments)

	return &SegmentUpdater{
		indexMeta:      indexMeta,
		index:          idx,
		segmentManager: segmentManager,
	}, nil
}
