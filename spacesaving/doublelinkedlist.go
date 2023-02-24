package spacesaving

type node struct {
	Bucket *bucket
	Prev   *node
	Next   *node
}

func newNode(b *bucket) *node {
	return &node{
		Bucket: b,
	}
}

type DoubleLinkedList struct {
	Head *node
	Tail *node
}

func newDoubleLinkedList() *DoubleLinkedList {
	return &DoubleLinkedList{}
}

func (ll *DoubleLinkedList) InsertAfter(n *node, newNode *node) {
	newNode.Prev = n
	newNode.Next = n.Next
	n.Next = newNode
	if ll.Tail == n {
		ll.Tail = newNode
	}
}

func (ll *DoubleLinkedList) InsertBefore(n *node, newNode *node) {
	newNode.Next = n
	newNode.Prev = n.Prev
	n.Prev = newNode
	if ll.Head == n {
		ll.Head = newNode
	}
}

func (ll *DoubleLinkedList) InsertBeginning(newNode *node) {
	if ll.Head == nil {
		ll.Head = newNode
		ll.Tail = newNode
	} else {
		ll.InsertBefore(ll.Head, newNode)
	}
}

func (ll *DoubleLinkedList) InsertEnd(newNode *node) {
	if ll.Tail == nil {
		ll.Head = newNode
		ll.Tail = newNode
	} else {
		ll.InsertAfter(ll.Tail, newNode)
	}
}

func (ll *DoubleLinkedList) Remove(n *node) {
	if n.Prev != nil {
		n.Prev.Next = n.Next
	} else {
		ll.Head = n.Next
	}
	if n.Next != nil {
		n.Next.Prev = n.Prev
	} else {
		ll.Tail = n.Prev
	}
}
