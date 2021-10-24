package index

import (
	"encoding/json"

	"github.com/k0kubun/pp/v3"

	"github.com/k-yomo/inverted-index/directory"
	"github.com/k-yomo/inverted-index/schema"
)

type Builder struct {
	schema *schema.Schema
}

func NewBuilder(indexSchema *schema.Schema) *Builder {
	return &Builder{schema: indexSchema}
}

func (b *Builder) CreateInDir(path string) (*Index, error) {
	mmapDirectory, err := directory.NewMMapDirectory(path)
	if err != nil {
		return nil, err
	}
	return b.create(mmapDirectory)
}

func (b *Builder) create(dir directory.Directory) (*Index, error) {
	managedDirectory, err := directory.NewManagedDirectory(dir)
	if err != nil {
		return nil, err
	}
	if err := b.saveNewMetas(managedDirectory); err != nil {
		return nil, err
	}
	indexMeta := NewIndexMeta(b.schema)
	return NewIndexFromMetas(managedDirectory, indexMeta, &SegmentMetaInventory{}), nil
}

func (b *Builder) saveNewMetas(dir directory.Directory) error {
	return b.saveMetas(NewIndexMeta(b.schema), dir)
}

func (b *Builder) saveMetas(meta *IndexMeta, dir directory.Directory) error {
	metaJSON, err := json.Marshal(meta)
	pp.Println(meta)
	if err != nil {
		return err
	}
	return dir.AtomicWrite(metaFileName, metaJSON)
}
