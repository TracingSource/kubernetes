// Copyright 2019, OpenCensus Authors
//


package trace

import (
	"github.com/hashicorp/golang-lru/simplelru"
)

type lruMap struct {
	simpleLruMap *simplelru.LRU
	droppedCount int
}

func newLruMap(size int) *lruMap {
	lm := &lruMap{}
	lm.simpleLruMap, _ = simplelru.NewLRU(size, nil)
	return lm
}

func (lm *lruMap) add(key, value interface{}) {
	evicted := lm.simpleLruMap.Add(key, value)
	if evicted {
		lm.droppedCount++
	}
}
