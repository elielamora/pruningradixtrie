package pruningradixtrie

type bsResults []Result

var _ ResultSet = &bsResults{}

func NewBSResultSet() *bsResults {
	return &bsResults{}
}

// NewBSResultSetFromSlice initializes a binary search result set
// with the contents of a result. This should only be used for testing
// and therefore it is the responibility of the caller to ensure that
// the results are sorted.
func NewBSResultSetFromSlice(results []Result) *bsResults {
	bs := NewBSResultSet()
	for _, r := range results {
		*bs = append(*bs, r)
	}
	return bs
}

func (bs bsResults) Len() int {
	return len(bs)
}

func (s bsResults) PeekMinResult() Result {
	return s[s.Len()-1]
}

func (bs *bsResults) PushResult(r Result) {
	bs.Insert(r, bs.Search(r.Freq))
}

// search finds the insertion index of the target
// implementation taken from stdlib sort/search.go
func (bs *bsResults) Search(freq uint64) int {
	arr := *bs
	i, j := 0, len(arr)
	for i < j {
		h := int(uint(i+j) >> 1)
		// i ≤ h < j
		if arr[h].Freq >= freq {
			i = h + 1
		} else {
			j = h
		}
	}
	return i
}

func (bs *bsResults) Insert(val Result, at int) {
	*bs = append(*bs, val)
	b := *bs
	// cannot swap with i-1 method for array of length 1
	// or if at is supposed to be appended
	if len(b) == at || len(b) < 2 {
		return
	}
	for i := len(b) - 1; i > at; i-- {
		b[i], b[i-1] = b[i-1], b[i]
	}
}

func (bs *bsResults) PopResult() Result {
	tail := bs.PeekMinResult()
	*bs = (*bs)[:bs.Len()-1]
	return tail
}

func (bs bsResults) Results() []Result {
	return bs
}
