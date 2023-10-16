package pruningradixtrie

import (
	"container/heap"
)

// rheap is a min-heap of Results
type rheap []Result

// resultHeap wraps a rheap
type resultHeap struct {
	h *rheap
}

var _ ResultSet = &resultHeap{}
var _ heap.Interface = &rheap{}

func NewResultHeap() ResultSet {
	return &resultHeap{
		h: new(rheap),
	}
}

func (r rheap) Len() int           { return len(r) }
func (r rheap) Less(i, j int) bool { return r[i].Freq < r[j].Freq }
func (r rheap) Swap(i, j int)      { r[i], r[j] = r[j], r[i] }

func (r *rheap) Push(x any) {
	*r = append(*r, x.(Result))
}

func (r rheap) PeekMin() any {
	return r[0]
}

func (r *rheap) Contents() []Result {
	return *r
}

func (r rheap) Slice(l int) []Result {
	return r[:l]
}

func (r *rheap) Pop() any {
	old := *r
	n := len(old)
	x := old[n-1]
	*r = old[0 : n-1]
	return x
}

func (r *resultHeap) PushResult(x Result) {
	heap.Push(r.h, x)
}

func (r *resultHeap) PeekMinResult() Result {
	return r.h.PeekMin().(Result)
}

func (r *resultHeap) PopResult() Result {
	return heap.Pop(r.h).(Result)
}

func (r *resultHeap) Len() int { return r.h.Len() }

// Results implements ResultSet.
// Returns results in sorted order,
// essentially performing heap sort in-place
// and returning the result.
func (r resultHeap) Results() []Result {
	h := *r.h
	// keep reference to full slice, will become the result once sorted
	results := h.Contents()
	for h.Len() > 0 {
		h.Swap(0, h.Len()-1)
		h = h.Slice(h.Len() - 1)
		heap.Fix(&h, 0)
	}
	return results
}
