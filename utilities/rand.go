package utilities

import (
	"math/rand"
	"sync"
	"time"
)

type rndSrc struct {
	mtx sync.Mutex
	src rand.Source
}

func (s *rndSrc) Int63() int64 {
	s.mtx.Lock()
	n := s.src.Int63()
	s.mtx.Unlock()
	return n
}

func (s *rndSrc) Seed(n int64) {
	s.mtx.Lock()
	s.src.Seed(n)
	s.mtx.Unlock()
}

func CreateRndSrc() *rand.Rand {
	return rand.New(&rndSrc{src: rand.NewSource(time.Now().UnixNano())})
}
