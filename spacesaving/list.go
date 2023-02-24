package spacesaving

type Node struct {
	Prev  *Node
	Next  *Node
	Value *Counter
}

type List struct {
	Head *Node
	Tail *Node
}

func newList() *List {
	return &List{
		Head: nil,
		Tail: nil,
	}
}

func (l *List) PushBack(v *Counter) *Node {
	n := &Node{
		Value: v,
	}
	if l.Head == nil {
		l.Head = n
		l.Tail = n
	} else {
		l.Tail.Next = n
		n.Prev = l.Tail
		l.Tail = n
	}
	return n
}

func (l *List) Front() *Node {
	return l.Head
}

func (l *List) Remove(c *Counter) {
	for n := l.Head; n != nil; n = n.Next {
		if n.Value == c {
			if n.Prev != nil {
				n.Prev.Next = n.Next
			} else {
				l.Head = n.Next
			}
			if n.Next != nil {
				n.Next.Prev = n.Prev
			} else {
				l.Tail = n.Prev
			}
			return
		}
	}
}

func (l *List) Len() int {
	count := 0
	for n := l.Head; n != nil; n = n.Next {
		count++
	}
	return count
}
