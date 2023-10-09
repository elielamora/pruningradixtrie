package pruningradixtrie

type node struct {
	key           string
	count         uint64
	maxChildCount uint64
	children      []*node
}

func newNode(term string, count uint64) *node {
	return &node{
		key:   term,
		count: count,
	}
}
