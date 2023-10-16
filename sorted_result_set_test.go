package pruningradixtrie_test

import (
	_ "embed"
	"testing"

	prtrie "github.com/elielamora/pruningradixtrie"
	"github.com/stretchr/testify/assert"
)

func TestSortedResultSetEmpty(t *testing.T) {
	s := prtrie.NewSortedResults()
	assert.Zero(t, s.Len())
}

func TestSortedResultsAdd(t *testing.T) {
	s := prtrie.NewSortedResults()
	s.PushResult(prtrie.Result{Term: "foo", Freq: 5})
	assert.Equal(t, 1, s.Len())
	assert.Equal(t, prtrie.Result{Term: "foo", Freq: 5}, s.PeekMinResult())

	s.PushResult(prtrie.Result{Term: "bar", Freq: 1})
	assert.Equal(t, 2, s.Len())
	assert.Equal(t, prtrie.Result{Term: "bar", Freq: 1}, s.PeekMinResult())

	s.PushResult(prtrie.Result{Term: "qux", Freq: 10})
	assert.Equal(t, 3, s.Len())
	assert.Equal(t, prtrie.Result{Term: "bar", Freq: 1}, s.PeekMinResult())

	s.PushResult(prtrie.Result{Term: "zip", Freq: 3})
	assert.Equal(t, 4, s.Len())
	assert.Equal(t, prtrie.Result{Term: "bar", Freq: 1}, s.PeekMinResult())

	assert.Equal(t, []prtrie.Result{
		{Term: "qux", Freq: 10},
		{Term: "foo", Freq: 5},
		{Term: "zip", Freq: 3},
		{Term: "bar", Freq: 1},
	}, s.Results())
}
