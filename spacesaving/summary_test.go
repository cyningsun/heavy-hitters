package spacesaving

import (
	"strings"
	"testing"
)

func TestXxx(t *testing.T) {
	sample := "a b b c c c e e e e e d d d d g g g g g g g f f f f f f"
	samples := strings.Split(sample, " ")
	summary := New(0.1)
	for i := 0; i < len(samples); i++ {
		summary.Offer(samples[i], 1)
	}
	want := newList()
	want.PushBack(&Counter{Value: 7, Item: "g"})
	want.PushBack(&Counter{Value: 6, Item: "f"})
	want.PushBack(&Counter{Value: 5, Item: "e"})
	p := want.Head
	topK := summary.TopK(3)
	for e := topK.Front(); e != nil; e = e.Next {
		if e.Value.Item != p.Value.Item || e.Value.Value != p.Value.Value {
			t.Fatalf("want %v, got %v", p.Value, e.Value)
		}
		p = p.Next
	}
}
