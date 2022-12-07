// Copyright 2019, OpenCensus Authors
//


package trace

type evictedQueue struct {
	queue        []interface{}
	capacity     int
	droppedCount int
}

func newEvictedQueue(capacity int) *evictedQueue {
	eq := &evictedQueue{
		capacity: capacity,
		queue:    make([]interface{}, 0),
	}

	return eq
}

func (eq *evictedQueue) add(value interface{}) {
	if len(eq.queue) == eq.capacity {
		eq.queue = eq.queue[1:]
		eq.droppedCount++
	}
	eq.queue = append(eq.queue, value)
}
