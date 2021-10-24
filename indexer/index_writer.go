package indexer

import (
	"fmt"
	"math"
	"runtime"

	"github.com/k-yomo/inverted-index/internal/opstamp"

	"github.com/k-yomo/inverted-index/index"
)

const (
	MaxThreadNum          = 8
	MarginInBytes         = 1e6 // 1MB
	HeapSizeMin           = MarginInBytes * 3
	HeapSizeMax           = math.MaxUint32 - MarginInBytes
	MaxOperationQueueSize = 10000
)

type IndexWriter struct {
	index             *index.Index
	heapSizePerThread int

	operationSender   chan<- []*AddOperation
	operationReceiver <-chan []*AddOperation
	segmentUpdater    *SegmentUpdater

	workerID  int
	ThreadNum int

	stamper          *opstamp.Stamper
	committedOpstamp opstamp.Opstamp
}

func NewIndexWriter(idx *index.Index, overallHeapBytes int) (*IndexWriter, error) {
	threadNum := int(math.Min(float64(runtime.GOMAXPROCS(0)), 8))
	heapSizePerThread := overallHeapBytes / threadNum
	if heapSizePerThread < HeapSizeMin {
		threadNum = int(math.Max(float64(overallHeapBytes/HeapSizeMin), 1))
	}
	if heapSizePerThread < HeapSizeMin {
		return nil, fmt.Errorf("heap size per thread needs to be at least %d", HeapSizeMin)
	}

	indexMeta, err := idx.LoadMetas()
	if err != nil {
		return nil, err
	}

	currentOpstamp := indexMeta.Opstamp

	stamper := opstamp.NewStamper(currentOpstamp)
	segmentUpdater, err := NewSegmentUpdater(idx, stamper)
	if err != nil {
		return nil, err
	}

	operationChan := make(chan []*AddOperation, MaxOperationQueueSize)

	i := &IndexWriter{
		index: idx,

		operationSender:   operationChan,
		operationReceiver: operationChan,
		segmentUpdater:    segmentUpdater,

		workerID:  0,
		ThreadNum: threadNum,

		stamper:          stamper,
		committedOpstamp: currentOpstamp,
	}

	if err := i.startWorkers(); err != nil {
		return nil, err
	}

	return i, nil
}

func (i *IndexWriter) startWorkers() error {
	for j := 0; j < i.ThreadNum; j++ {
		if err := i.addIndexWorker(); err != nil {
			return err
		}
	}
	return nil
}

func (i *IndexWriter) addIndexWorker() error {
	go func() {
		for {
			operations := <-i.operationReceiver
			if err := i.indexDocuments(operations); err != nil {
				// logging?
			}
		}
	}()

	i.workerID += 1
	return nil
}

func (i *IndexWriter) indexDocuments(operations []*AddOperation) error {
	segment := i.index.NewSegment()
	schema := segment.Schema()
	segmentWriter := newSegmentWriter(i.heapSizePerThread, segment, schema)
	for _, op := range operations {
		if err := segmentWriter.addDocument(op, schema); err != nil {
			return err
		}
	}
}
