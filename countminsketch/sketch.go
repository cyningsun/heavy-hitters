package countminsketch

import (
	"container/heap"
	"hash/fnv"
	"log"
	"math"

	"github.com/cyningsun/heavy-hitters/internal/queue"
)

// https://zhuanlan.zhihu.com/p/369981005

type summary struct {
	sketch *CountMinSketch
	k      int
	queue  queue.PriorityQueue
	counts map[string]*queue.Item
}

func New(sketch *CountMinSketch, k int) *summary {
	return &summary{
		sketch: sketch,
		k:      k,
		queue:  make(queue.PriorityQueue, 0, k),
		counts: make(map[string]*queue.Item, k),
	}
}

func (s *summary) Incr(str string) {
	count := s.sketch.Increment(str)
	_, ok := s.counts[str]
	switch {
	case !ok && len(s.queue) < s.k:
		it := &queue.Item{
			Value: str,
			Count: count,
			Index: 0,
		}

		heap.Push(&s.queue, it)
		s.counts[str] = it
	case !ok && len(s.queue) >= s.k && count > s.queue[0].Count:
		minItem := heap.Pop(&s.queue).(*queue.Item)
		delete(s.counts, minItem.Value)
		minItem.Value = str
		minItem.Count = count
		s.counts[str] = minItem
		heap.Push(&s.queue, minItem)
	case !ok && len(s.queue) >= s.k && count <= s.queue[0].Count: // do nothing
	case ok:
		s.counts[str].Count = count
		heap.Fix(&s.queue, s.counts[str].Index)
	}
}

func (s *summary) TopK() map[string]int {
	result := make(map[string]int)
	for _, each := range s.queue {
		result[each.Value] = each.Count
	}

	return result
}

// WithEstimates creates a Count-Min Sketch with given error rate and confidence.
// Accuracy guarantees will be made in terms of a pair of user specified parameters,
// ε and δ, meaning that the error in answering a query is within a factor of ε with
// probability δ
func WithEstimates(epsilon, delta float64) (sk *CountMinSketch) {
	const (
		defaultErrorRatio  = 1.0 / 1e3 // 0.1%
		defaultUncertainty = 1.0 / 1e3 // 0.1%
	)
	if epsilon < 1.0/1e9 || epsilon > 0.1 {
		log.Printf("error ratio %g not in [1e-9,0.1], use default %g\n", epsilon, defaultErrorRatio)
		epsilon = defaultErrorRatio
	}
	if delta < 1.0/1e9 || delta > 0.1 {
		log.Printf("certainty %g not in [1e-9,0.1], use default %g\n", delta, defaultUncertainty)
		delta = defaultUncertainty
	}

	width := uint32(math.Ceil(2 / epsilon))
	depth := uint32(math.Ceil(math.Log(1/delta) / math.Log(2)))
	return NewCountMinSketch(depth, width)
}

type CountMinSketch struct {
	numHashes  uint32
	numBuckets uint32
	hashFuncs  []func(item string) uint32
	counts     [][]int
}

func NewCountMinSketch(numHashes, numBuckets uint32) *CountMinSketch {
	hashFuncs := make([]func(item string) uint32, numHashes)
	for i := range hashFuncs {
		h := fnv.New32a()
		h.Write([]byte{byte(i)})
		hashFuncs[i] = func(item string) uint32 {
			h.Reset()
			h.Write([]byte(item))
			return h.Sum32()
		}
	}

	counts := make([][]int, numHashes)
	for i := range counts {
		counts[i] = make([]int, numBuckets)
	}

	return &CountMinSketch{
		numHashes:  numHashes,
		hashFuncs:  hashFuncs,
		counts:     counts,
		numBuckets: numBuckets,
	}
}

func (cms *CountMinSketch) Increment(item string) int {
	minCount := int(^uint(0) >> 1)
	for i, hashFunc := range cms.hashFuncs {
		hashValue := hashFunc(item)
		bucket := hashValue % uint32(cms.numBuckets)
		cms.counts[i][bucket]++

		if cms.counts[i][bucket] < minCount {
			minCount = cms.counts[i][bucket]
		}
	}
	return minCount
}

func (cms *CountMinSketch) Estimate(item string) int {
	minCount := int(^uint(0) >> 1)
	for i, hashFunc := range cms.hashFuncs {
		hashValue := hashFunc(item)
		bucket := hashValue % uint32(cms.numBuckets)
		count := cms.counts[i][bucket]
		if count < minCount {
			minCount = count
		}
	}
	return minCount
}
