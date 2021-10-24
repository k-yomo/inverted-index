package opstamp

import (
	"sync"
)

type Opstamp uint64

type Stamper struct {
	mu      sync.Mutex
	opstamp Opstamp
}

func NewStamper(firstOpstamp Opstamp) *Stamper {
	return &Stamper{opstamp: firstOpstamp}
}

func (s *Stamper) Stamp() Opstamp {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.opstamp += 1
	return s.opstamp
}
