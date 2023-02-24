package spacesaving

type Counter struct {
	Value int
	Error int
	Item  string
	Node  *node
}

func newCounter(node *node) *Counter {
	return &Counter{
		Value: 0,
		Error: 0,
		Node:  node,
	}
}
