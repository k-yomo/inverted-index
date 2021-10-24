package indexer

import (
	"github.com/k-yomo/inverted-index/internal/opstamp"
	"github.com/k-yomo/inverted-index/schema"
)

type AddOperation struct {
	opstamp  opstamp.Opstamp
	document *schema.Document
}
