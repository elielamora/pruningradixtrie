package pruningradixtrie

import (
	"sort"
)

type sortedResults []Result

var _ ResultSet = &sortedResults{}

func NewSortedResults() *sortedResults {
	return &sortedResults{}
}

func (s sortedResults) Len() int {
	return len(s)
}

func (s sortedResults) PeekMinResult() Result {
	return s[s.Len()-1]
}

func (s *sortedResults) PushResult(r Result) {
	*s = append(*s, r)
	sr := *s
	sort.Slice(sr, func(i, j int) bool {
		return sr[j].Freq < sr[i].Freq
	})
}

// PopResult implements ResultSet.
func (s *sortedResults) PopResult() Result {
	tail := s.PeekMinResult()
	*s = (*s)[:s.Len()-1]
	return tail
}

func (s sortedResults) Results() []Result {
	return s
}
