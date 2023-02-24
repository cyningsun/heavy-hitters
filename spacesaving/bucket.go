package spacesaving

type bucket struct {
	Children *List
	Value    int
}

func newBucket(value int) *bucket {
	return &bucket{
		Children: newList(),
		Value:    value,
	}
}
