package spacesaving

import (
	"container/list"
	"fmt"
	"math"
	"sync/atomic"
)

type summary struct {
	Capacity int
	Count    atomic.Uint64
	Cache    map[string]*Counter
	Buckets  *DoubleLinkedList
}

func New(epsilon float64) *summary {
	const (
		defaultErrorRatio = 1.0 / 1e3 // 0.1%
	)
	if epsilon < 1.0/1e9 || epsilon > 0.1 {
		fmt.Printf("error ratio %g not in [1e-9,0.1], use default %g\n", epsilon, defaultErrorRatio)
		epsilon = defaultErrorRatio
	}
	capacity := int(math.Ceil(1.0 / epsilon))
	newBucket := newBucket(0)

	node := newNode(newBucket)
	for i := 0; i < capacity; i++ {
		newBucket.Children.PushBack(newCounter(node))
	}

	buckets := newDoubleLinkedList()
	buckets.InsertBeginning(node)

	return &summary{
		Capacity: capacity,
		Count:    atomic.Uint64{},
		Cache:    make(map[string]*Counter),
		Buckets:  buckets,
	}
}

func (s *summary) TopK(k int) *list.List {
	head := list.New()
	bucket := s.Buckets.Head
	for bucket != nil {
		child := bucket.Bucket.Children.Front()
		for child != nil {
			head.PushBack(child.Value)
			child = child.Next
		}
		bucket = bucket.Next
		if head.Len() == k {
			break
		}
	}
	return head
}

func (s *summary) incrementCounter(counter *Counter) {
	node := counter.Node
	bucketNext := node.Prev
	node.Bucket.Children.Remove(counter)
	counter.Value = counter.Value + 1

	if bucketNext != nil && counter.Value == bucketNext.Bucket.Value {
		bucketNext.Bucket.Children.PushBack(counter)
		counter.Node = bucketNext
	} else {
		newBucket := newBucket(counter.Value)
		newBucket.Children.PushBack(counter)
		s.Buckets.InsertBefore(node, newNode(newBucket))
		counter.Node = node.Prev
	}

	if node.Bucket.Children.Len() == 0 {
		s.Buckets.Remove(node)
	}
}

func (s *summary) Offer(item string, increment int) {
	s.Count.Add(1)

	if _, ok := s.Cache[item]; ok {
		counter := s.Cache[item]
		s.incrementCounter(counter)
		return
	}

	minElement := s.Buckets.Tail.Bucket.Children.Front().Value
	originalMinValue := minElement.Value
	delete(s.Cache, minElement.Item)
	s.Cache[item] = minElement
	minElement.Item = item
	s.incrementCounter(minElement)
	if len(s.Cache) <= s.Capacity {
		minElement.Error = originalMinValue
	}
}
