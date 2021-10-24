package index

import (
	"github.com/k-yomo/inverted-index/internal/opstamp"
	"github.com/k-yomo/inverted-index/schema"
)

const metaFileName = "meta.json"

type IndexMeta struct {
	Segments []*SegmentMeta `json:"segments"`
	Schema   *schema.Schema `json:"schema"`
	// last commit operation's id
	Opstamp opstamp.Opstamp `json:"opstamp"`
}

type SegmentMeta struct {
	SegmentID SegmentID   `json:"segmentId"`
	MacDoc    uint32      `json:"maxDoc"`
	Deletes   *DeleteMeta `json:"deletes"`
}

type SegmentMetaInventory struct {
	inventory []*SegmentMeta
}

type DeleteMeta struct {
	NumDeletedDocs int             `json:"numDeletedDocs"`
	Opstamp        opstamp.Opstamp `json:"operationId"`
}

func NewIndexMeta(schema *schema.Schema) *IndexMeta {
	return &IndexMeta{
		Segments: nil,
		Schema:   schema,
		Opstamp:  0,
	}
}

func (i *SegmentMetaInventory) NewSegmentMeta(segmentID SegmentID, maxDoc uint32) *SegmentMeta {
	segmentMeta := &SegmentMeta{
		SegmentID: segmentID,
		MacDoc:    maxDoc,
		Deletes:   nil,
	}
	// TODO: Make it thread safe
	i.inventory = append(i.inventory, segmentMeta)
	return segmentMeta
}
