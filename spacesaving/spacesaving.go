package main

import (
	"container/heap"
	"fmt"
	"math"

	"github.com/cyningsun/heavy-hitters/internal/queue"
)

type summary struct {
	queue      queue.PriorityQueue
	counts     map[string]*queue.Item
	counterNum int
	k          int
}

// New returns a new summary with the given error ratio and k.
func New(epsilon float64, k int) *summary {
	const (
		defaultErrorRatio = 1.0 / 1e3 // 0.1%
	)
	if epsilon < 1.0/1e9 || epsilon > 0.1 {
		fmt.Printf("error ratio %g not in [1e-9,0.1], use default %g\n", epsilon, defaultErrorRatio)
		epsilon = defaultErrorRatio
	}

	num := int(math.Ceil(1.0 / epsilon))

	return &summary{
		queue:      make(queue.PriorityQueue, 0, num),
		counts:     make(map[string]*queue.Item, num),
		counterNum: num,
		k:          k,
	}
}

func (s *summary) Incr(str string) {
	_, ok := s.counts[str]
	switch {
	case !ok && len(s.queue) < s.counterNum:
		it := &queue.Item{
			Value:    str,
			Count:    1,
			ErrCount: 0,
			Index:    0,
		}

		heap.Push(&s.queue, it)
		s.counts[str] = it
	case !ok && len(s.queue) >= s.counterNum:
		minItem := heap.Pop(&s.queue).(*queue.Item)
		delete(s.counts, minItem.Value)

		minItem.Value = str
		minItem.ErrCount = minItem.Count
		minItem.Count++
		heap.Push(&s.queue, minItem)
		s.counts[str] = minItem
	case ok:
		s.counts[str].Count++
		heap.Fix(&s.queue, s.counts[str].Index)
	}
}

func (s *summary) TopK() map[string]int {
	for len(s.queue) > s.k {
		heap.Pop(&s.queue)
	}

	result := make(map[string]int)
	for _, each := range s.queue {
		result[each.Value] = each.Count
	}

	return result
}
