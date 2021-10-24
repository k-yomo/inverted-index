package index

import (
	"encoding/json"
	"sync"

	"github.com/k-yomo/inverted-index/analyzer"
	"github.com/k-yomo/inverted-index/directory"
	"github.com/k-yomo/inverted-index/schema"
)

type Index struct {
	directory directory.Directory
	schema    *schema.Schema
	inventory *SegmentMetaInventory

	Analyzer analyzer.Analyzer

	mu *sync.Mutex
}

func NewIndexFromMetas(directory directory.Directory, metas *IndexMeta, inventory *SegmentMetaInventory) *Index {
	return &Index{
		directory: directory,
		schema:    metas.Schema,
		inventory: inventory,
		Analyzer:  &analyzer.EnglishAnalyzer{},
		mu:        &sync.Mutex{},
	}
}

func (i *Index) LoadMetas() (*IndexMeta, error) {
	metaData, err := i.directory.AtomicRead(metaFileName)
	if err != nil {
		return nil, err
	}
	var indexMeta IndexMeta
	if err := json.Unmarshal(metaData, &indexMeta); err != nil {
		return nil, err
	}

	return &indexMeta, nil
}

func (i *Index) NewSegment() *Segment {
	segmentMeta := i.inventory.NewSegmentMeta(NewSegmentID(), 0)
	return newSegment(i, segmentMeta)
}
