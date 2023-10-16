package pruningradixtrie_test

import (
	_ "embed"
	"testing"

	prtrie "github.com/elielamora/pruningradixtrie"
	"github.com/stretchr/testify/assert"
)

func TestHeapEmpty(t *testing.T) {
	h := prtrie.NewResultHeap()
	assert.Zero(t, h.Len())
}

func TestHeapAdd(t *testing.T) {
	h := prtrie.NewResultHeap()
	h.PushResult(prtrie.Result{Term: "foo", Freq: 5})
	assert.Equal(t, 1, h.Len())
	assert.Equal(t, prtrie.Result{Term: "foo", Freq: 5}, h.PeekMinResult())

	h.PushResult(prtrie.Result{Term: "bar", Freq: 1})
	assert.Equal(t, 2, h.Len())
	assert.Equal(t, prtrie.Result{Term: "bar", Freq: 1}, h.PeekMinResult())

	h.PushResult(prtrie.Result{Term: "qux", Freq: 10})
	assert.Equal(t, 3, h.Len())
	assert.Equal(t, prtrie.Result{Term: "bar", Freq: 1}, h.PeekMinResult())

	h.PushResult(prtrie.Result{Term: "zip", Freq: 3})
	assert.Equal(t, 4, h.Len())
	assert.Equal(t, prtrie.Result{Term: "bar", Freq: 1}, h.PeekMinResult())

	assert.Equal(t, []prtrie.Result{
		{Term: "qux", Freq: 10},
		{Term: "foo", Freq: 5},
		{Term: "zip", Freq: 3},
		{Term: "bar", Freq: 1},
	}, h.Results())
}
