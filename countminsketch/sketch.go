package countminsketch

import (
	"encoding/binary"
	"fmt"
	"hash"
	"hash/fnv"
	"math"
	"sync"
	"unsafe"
)

type CountType uint32

const Max = ^(CountType(0))

type Sketch struct {
	width  uint32
	depth  uint32
	count  [][]CountType
	hasher hash.Hash64
	mutex  sync.RWMutex
}

// WithEstimates creates a Count-Min Sketch with given error rate and confidence.
// Accuracy guarantees will be made in terms of a pair of user specified parameters,
// ε and δ, meaning that the error in answering a query is within a factor of ε with
// probability δ
func WithEstimates(epsilon, delta float64) (sk *Sketch) {
	const (
		defaultErrorRatio  = 1.0 / 1e3 // 0.1%
		defaultUncertainty = 1.0 / 1e3 // 0.1%
	)
	if epsilon < 1.0/1e9 || epsilon > 0.1 {
		fmt.Printf("error ratio %g not in [1e-9,0.1], use default %g\n", epsilon, defaultErrorRatio)
		epsilon = defaultErrorRatio
	}
	if delta < 1.0/1e9 || delta > 0.1 {
		fmt.Printf("certainty %g not in [1e-9,0.1], use default %g\n", delta, defaultUncertainty)
		delta = defaultUncertainty
	}

	width := uint32(math.Ceil(2 / epsilon))
	depth := uint32(math.Ceil(-math.Log(delta) / math.Log(2)))
	return New(width, depth)
}

// New returns a new Count-Min Sketch with the given width and depth.
func New(width, depth uint32) (sk *Sketch) {
	sk = &Sketch{
		width:  width,
		depth:  depth,
		count:  make([][]CountType, depth),
		hasher: fnv.New64(),
	}
	for i := uint32(0); i < depth; i++ {
		sk.count[i] = make([]CountType, width)
	}
	return sk
}

// Width returns the width of the sketch.
func (sk *Sketch) Width() uint32 { return sk.width }

// Depth returns the depth of the sketch.
func (sk *Sketch) Depth() uint32 { return sk.depth }

// String returns a string representation of the sketch.
func (sk *Sketch) String() string {
	space := float64(int64(sk.width)*int64(sk.depth)*int64(unsafe.Sizeof(sk.count[0][0]))) / 1e6
	return fmt.Sprintf("Count-Min Sketch(%p): width=%d, depth=%d, mem=%.3fm",
		sk, sk.width, sk.depth, space)
}

// Clear resets the sketch.
func (sk *Sketch) Clear() {
	sk.mutex.Lock()
	for i := uint32(0); i < sk.depth; i++ {
		for j := uint32(0); j < sk.width; j++ {
			sk.count[i][j] = 0
		}
	}
	sk.mutex.Unlock()
}

// Incr increments the count for the given key by 1.
func (sk *Sketch) Incr(key []byte) (min CountType) {
	return sk.Add(key, 1)
}

// Add adds cnt to the count for the given key.
func (sk *Sketch) Add(key []byte, cnt CountType) (min CountType) {
	pos := sk.positions(key)
	min = sk.query(pos)

	min += cnt

	sk.mutex.Lock()
	for i := uint32(0); i < sk.depth; i++ {
		v := sk.count[i][pos[i]]
		if v < min {
			sk.count[i][pos[i]] = min
		}
	}
	sk.mutex.Unlock()

	return min
}

// Query returns the minimum count for the given key.
// If the key does not exist, returns 0.
func (sk *Sketch) Query(key []byte) (min CountType) {
	pos := sk.positions(key)
	return sk.query(pos)
}

// hashs returns the two hashes of key.
// It uses the FNV hash function to compute the hashes.
func (s *Sketch) hashs(key []byte) (a uint32, b uint32) {
	s.hasher.Reset()
	s.hasher.Write(key)
	sum := s.hasher.Sum(nil)
	upper := sum[0:4]
	lower := sum[4:8]
	a = binary.BigEndian.Uint32(lower)
	b = binary.BigEndian.Uint32(upper)
	return
}

// positions returns the indices of the counters that should be updated for a
// given key. The indices are calculated using the double hashing technique.
func (sk *Sketch) positions(key []byte) (pos []uint32) {
	hash1, hash2 := sk.hashs(key)
	pos = make([]uint32, sk.depth)
	for i := uint32(0); i < sk.depth; i++ {
		pos[i] = (hash1 + i*hash2) % sk.width
	}
	return pos
}

// query returns the minimum count for the given positions.
func (sk *Sketch) query(pos []uint32) (min CountType) {
	min = Max

	sk.mutex.RLock()
	for i := uint32(0); i < sk.depth; i++ {
		v := sk.count[i][pos[i]]
		if min > v {
			min = v
		}
	}
	sk.mutex.RUnlock()

	return min
}
