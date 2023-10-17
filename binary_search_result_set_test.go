package pruningradixtrie_test

import (
	_ "embed"
	"testing"

	prtrie "github.com/elielamora/pruningradixtrie"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBSResultSetEmpty(t *testing.T) {
	bs := prtrie.NewBSResultSet()
	assert.Zero(t, bs.Len())
}

func TestBSResultSetAdd(t *testing.T) {
	bs := prtrie.NewBSResultSet()
	bs.PushResult(prtrie.Result{Term: "foo", Freq: 5})
	assert.Equal(t, 1, bs.Len())
	assert.Equal(t, prtrie.Result{Term: "foo", Freq: 5}, bs.PeekMinResult())

	bs.PushResult(prtrie.Result{Term: "bar", Freq: 1})
	assert.Equal(t, 2, bs.Len())
	assert.Equal(t, prtrie.Result{Term: "bar", Freq: 1}, bs.PeekMinResult())

	bs.PushResult(prtrie.Result{Term: "qux", Freq: 10})
	assert.Equal(t, 3, bs.Len())
	assert.Equal(t, prtrie.Result{Term: "bar", Freq: 1}, bs.PeekMinResult())

	bs.PushResult(prtrie.Result{Term: "zip", Freq: 3})
	assert.Equal(t, 4, bs.Len())
	assert.Equal(t, prtrie.Result{Term: "bar", Freq: 1}, bs.PeekMinResult())

	assert.Equal(t, []prtrie.Result{
		{Term: "qux", Freq: 10},
		{Term: "foo", Freq: 5},
		{Term: "zip", Freq: 3},
		{Term: "bar", Freq: 1},
	}, bs.Results())
}

func TestSearch(t *testing.T) {
	bs := prtrie.NewBSResultSetFromSlice([]prtrie.Result{
		{Term: "0", Freq: 100},
		{Term: "1", Freq: 75},
		{Term: "2", Freq: 50},
		{Term: "3", Freq: 25},
		{Term: "4", Freq: 10},
		{Term: "5", Freq: 5},
	})
	assert.Equal(t, 6, bs.Len())

	assert.Equal(t, 0, bs.Search(99999))
	assert.Equal(t, 1, bs.Search(100))
	assert.Equal(t, 1, bs.Search(99))
	assert.Equal(t, 2, bs.Search(75))
	assert.Equal(t, 2, bs.Search(70))
	assert.Equal(t, 3, bs.Search(50))
	assert.Equal(t, 3, bs.Search(49))
	assert.Equal(t, 4, bs.Search(25))
	assert.Equal(t, 4, bs.Search(20))
	assert.Equal(t, 5, bs.Search(10))
	assert.Equal(t, 5, bs.Search(8))
	assert.Equal(t, 6, bs.Search(5))
	assert.Equal(t, 6, bs.Search(1))
}

func TestInsert(t *testing.T) {
	bs := prtrie.NewBSResultSet()

	bs.Insert(prtrie.Result{Term: "2", Freq: 50}, 4)
	require.Equal(t, 1, bs.Len())
	assert.Equal(t, []prtrie.Result{
		// {Term: "0", Freq: 100},
		// {Term: "1", Freq: 75},
		{Term: "2", Freq: 50},
		// {Term: "3", Freq: 25},
		// {Term: "4", Freq: 10},
		// {Term: "5", Freq: 5},
	}, bs.Results())

	bs.Insert(prtrie.Result{Term: "1", Freq: 75}, 0)
	require.Equal(t, 2, bs.Len())
	assert.Equal(t, []prtrie.Result{
		// {Term: "0", Freq: 100},
		{Term: "1", Freq: 75},
		{Term: "2", Freq: 50},
		// {Term: "3", Freq: 25},
		// {Term: "4", Freq: 10},
		// {Term: "5", Freq: 5},
	}, bs.Results())

	bs.Insert(prtrie.Result{Term: "4", Freq: 10}, 2)
	require.Equal(t, 3, bs.Len())
	assert.Equal(t, []prtrie.Result{
		// {Term: "0", Freq: 100},
		{Term: "1", Freq: 75},
		{Term: "2", Freq: 50},
		// {Term: "3", Freq: 25},
		{Term: "4", Freq: 10},
		// {Term: "5", Freq: 5},
	}, bs.Results())

	bs.Insert(prtrie.Result{Term: "5", Freq: 5}, 4)
	require.Equal(t, 4, bs.Len())
	assert.Equal(t, []prtrie.Result{
		// {Term: "0", Freq: 100},
		{Term: "1", Freq: 75},
		{Term: "2", Freq: 50},
		// {Term: "3", Freq: 25},
		{Term: "4", Freq: 10},
		{Term: "5", Freq: 5},
	}, bs.Results())

	bs.Insert(prtrie.Result{Term: "0", Freq: 100}, 0)
	require.Equal(t, 5, bs.Len())
	assert.Equal(t, []prtrie.Result{
		{Term: "0", Freq: 100},
		{Term: "1", Freq: 75},
		{Term: "2", Freq: 50},
		// {Term: "3", Freq: 25},
		{Term: "4", Freq: 10},
		{Term: "5", Freq: 5},
	}, bs.Results())

	bs.Insert(prtrie.Result{Term: "3", Freq: 25}, 3)
	require.Equal(t, 6, bs.Len())
	assert.Equal(t, []prtrie.Result{
		{Term: "0", Freq: 100},
		{Term: "1", Freq: 75},
		{Term: "2", Freq: 50},
		{Term: "3", Freq: 25},
		{Term: "4", Freq: 10},
		{Term: "5", Freq: 5},
	}, bs.Results())
}

func TestPopResult(t *testing.T) {
	bs := prtrie.NewBSResultSetFromSlice([]prtrie.Result{
		{Term: "0", Freq: 100},
		{Term: "1", Freq: 75},
		{Term: "2", Freq: 50},
		{Term: "3", Freq: 25},
		{Term: "4", Freq: 10},
		{Term: "5", Freq: 5},
	})
	assert.Equal(t, 6, bs.Len())
	popped := bs.PopResult()
	assert.Equal(t, 5, bs.Len())
	assert.Equal(t, prtrie.Result{Term: "5", Freq: 5}, popped)
}
