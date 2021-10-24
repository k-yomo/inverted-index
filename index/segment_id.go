package index

import "github.com/k-yomo/inverted-index/pkg/uuid"

type SegmentID string

func NewSegmentID() SegmentID {
	return SegmentID(uuid.Generate())
}
