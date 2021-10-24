package indexer

import (
	"sync"

	"github.com/k-yomo/inverted-index/index"
)

type SegmentManager struct {
	mu        sync.Mutex
	registers *SegmentRegisters
}

type SegmentRegisters struct {
	uncommitted *SegmentRegister
	committed   *SegmentRegister
}

func NewSegmentManager(segmentMetas []*index.SegmentMeta) *SegmentManager {
	return &SegmentManager{
		registers: &SegmentRegisters{
			uncommitted: newSegmentRegister(),
			committed:   newSegmentRegisterFromSegmentMetas(segmentMetas),
		},
	}
}
