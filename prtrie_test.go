package pruningradixtrie_test

import (
	"strings"
	"testing"

	prtrie "github.com/elielamora/pruningradixtrie"
	"github.com/stretchr/testify/assert"
)

func TestPrtrieEmpty(t *testing.T) {
	p := prtrie.NewPruningRadixTrie()
	assert.Equal(t, uint64(0), p.GetTotalTermCount(), "expected count to be zero")
	assert.Equal(t, 0, len(p.TopKForPrefix("", 1)))
}

func TestPrtrieAddEmptyString(t *testing.T) {
	p := prtrie.NewPruningRadixTrie()
	p.AddTerm("", 0)
	assert.Equal(t, uint64(0), p.GetTotalTermCount(), "expected count to be zero")
	assert.Equal(t, 0, len(p.TopKForPrefix("", 1)))

	p.AddTerm("", 1)
	assert.Equal(t, uint64(0), p.GetTotalTermCount(), "expected count to be zero")
	assert.Equal(t, 0, len(p.TopKForPrefix("", 1)))
}

func TestPrtrieAddTermWithZeroCount(t *testing.T) {
	p := prtrie.NewPruningRadixTrie()
	p.AddTerm("test", 0)
	assert.Equal(t, uint64(0), p.GetTotalTermCount(), "expected count to be zero")
	assert.Empty(t, p.TopKForPrefix("", 1))
}

func TestPrtrieAddSingle(t *testing.T) {
	p := prtrie.NewPruningRadixTrie()
	p.AddTerm("test", 7)
	assert.Equal(t, uint64(1), p.GetTotalTermCount(), "expected count to be one")
	assert.Equal(t, 0, len(p.TopKForPrefix("", 0)), "expect no top k results when k is 0")
	assert.Equal(t, 0, len(p.TopKForPrefix("tests", 1)), "expect no top k results when prefix shares start with term but is longer")
	assert.Equal(t, 0, len(p.TopKForPrefix("x", 1)), "expect no top k results when prefix is different")
	assert.Equal(t, []prtrie.Result{{Term: "test", Freq: 7}}, p.TopKForPrefix("", 1), "expected term added as the only result")
	assert.Equal(t, []prtrie.Result{{Term: "test", Freq: 7}}, p.TopKForPrefix("t", 1), "expected term added as the only result")
	assert.Equal(t, []prtrie.Result{{Term: "test", Freq: 7}}, p.TopKForPrefix("te", 1), "expected term added as the only result")
	assert.Equal(t, []prtrie.Result{{Term: "test", Freq: 7}}, p.TopKForPrefix("tes", 1), "expected term added as the only result")
	assert.Equal(t, []prtrie.Result{{Term: "test", Freq: 7}}, p.TopKForPrefix("test", 1), "expected term added as the only result")
	assert.Equal(t, []prtrie.Result{{Term: "test", Freq: 7}}, p.TopKForPrefix("", 2), "expected term added as the only result when k is greater than total terms")
}

func TestPrtrieAddTermMultipleTimes(t *testing.T) {
	p := prtrie.NewPruningRadixTrie()
	p.AddTerm("test", 6)
	p.AddTerm("test", 1)
	assert.Equal(t, uint64(1), p.GetTotalTermCount(), "expected count to be one")
	assert.Equal(t, 0, len(p.TopKForPrefix("", 0)), "expect no top k results when k is 0")
	assert.Equal(t, 0, len(p.TopKForPrefix("tests", 1)), "expect no top k results when prefix shares start with term but is longer")
	assert.Equal(t, 0, len(p.TopKForPrefix("x", 1)), "expect no top k results when prefix is different")
	assert.Equal(t, []prtrie.Result{{Term: "test", Freq: 7}}, p.TopKForPrefix("", 1), "expected term added as the only result")
	assert.Equal(t, []prtrie.Result{{Term: "test", Freq: 7}}, p.TopKForPrefix("t", 1), "expected term added as the only result")
	assert.Equal(t, []prtrie.Result{{Term: "test", Freq: 7}}, p.TopKForPrefix("te", 1), "expected term added as the only result")
	assert.Equal(t, []prtrie.Result{{Term: "test", Freq: 7}}, p.TopKForPrefix("tes", 1), "expected term added as the only result")
	assert.Equal(t, []prtrie.Result{{Term: "test", Freq: 7}}, p.TopKForPrefix("test", 1), "expected term added as the only result")
	assert.Equal(t, []prtrie.Result{{Term: "test", Freq: 7}}, p.TopKForPrefix("", 2), "expected term added as the only result when k is greater than total terms")
}

func TestPrtrieAddDisjointTerms(t *testing.T) {
	p := prtrie.NewPruningRadixTrie()
	p.AddTerm("foo", 7)
	p.AddTerm("bar", 77)
	p.AddTerm("qux", 777)
	assert.Equal(t, uint64(3), p.GetTotalTermCount(), "expected count to be 3")
	assert.Equal(t, 0, len(p.TopKForPrefix("", 0)), "expect no top k results when k is 0")
	assert.Equal(t, []prtrie.Result{{Term: "qux", Freq: 777}}, p.TopKForPrefix("", 1), "expected top term for empty prefix")
	assert.Equal(t, []prtrie.Result{
		{Term: "qux", Freq: 777},
		{Term: "bar", Freq: 77},
	}, p.TopKForPrefix("", 2), "expected top terms for empty prefix")
	assert.Equal(t, []prtrie.Result{
		{Term: "qux", Freq: 777},
		{Term: "bar", Freq: 77},
		{Term: "foo", Freq: 7},
	}, p.TopKForPrefix("", 3), "expected top terms for empty prefix")
	assert.Equal(t, []prtrie.Result{
		{Term: "qux", Freq: 777},
		{Term: "bar", Freq: 77},
		{Term: "foo", Freq: 7},
	}, p.TopKForPrefix("", 4), "expected top terms for empty prefix with k larger than results")
	assert.Equal(t, []prtrie.Result{
		{Term: "qux", Freq: 777},
	}, p.TopKForPrefix("q", 3), "expected top terms for prefix 'q'")
	assert.Equal(t, []prtrie.Result{
		{Term: "bar", Freq: 77},
	}, p.TopKForPrefix("b", 3), "expected top terms for prefix 'b'")
	assert.Equal(t, []prtrie.Result{
		{Term: "foo", Freq: 7},
	}, p.TopKForPrefix("f", 3), "expected top terms for prefix 'f'")
}

func TestPrtrieAddingSuffix(t *testing.T) {
	p := prtrie.NewPruningRadixTrie()
	p.AddTerm("test", 2)
	p.AddTerm("testing", 3)
	p.AddTerm("testing things", 1)
	assert.Equal(t, []prtrie.Result{
		{Term: "testing", Freq: 3},
		{Term: "test", Freq: 2},
		{Term: "testing things", Freq: 1},
	}, p.TopKForPrefix("", 3), "expected all top terms")
	assert.Equal(t, []prtrie.Result{
		{Term: "testing", Freq: 3},
		{Term: "test", Freq: 2},
	}, p.TopKForPrefix("test", 2), "expected top terms with the common prefix")
	assert.Equal(t, []prtrie.Result{
		{Term: "testing", Freq: 3},
		{Term: "testing things", Freq: 1},
	}, p.TopKForPrefix("testi", 3), "expected top terms with the common prefix")
	assert.Equal(t, []prtrie.Result{
		{Term: "testing things", Freq: 1},
	}, p.TopKForPrefix("testing th", 2), "expected top terms with the common prefix")
}

func TestPrtrieAddPrefixTerm(t *testing.T) {
	p := prtrie.NewPruningRadixTrie()
	p.AddTerm("foosball", 3)
	p.AddTerm("foo", 2)
	assert.Equal(t, []prtrie.Result{
		{Term: "foosball", Freq: 3},
		{Term: "foo", Freq: 2},
	}, p.TopKForPrefix("", 2), "expected all top terms")
	p.AddTerm("foo", 7)
	assert.Equal(t, []prtrie.Result{
		{Term: "foo", Freq: 9},
		{Term: "foosball", Freq: 3},
	}, p.TopKForPrefix("foo", 2), "expected all top terms")
	assert.Equal(t, []prtrie.Result{
		{Term: "foosball", Freq: 3},
	}, p.TopKForPrefix("foos", 2), "expected all top terms")
}

func TestPrtrieAddingOverlappingTerm(t *testing.T) {
	p := prtrie.NewPruningRadixTrie()
	p.AddTerm("testing", 5)
	p.AddTerm("tester", 10)
	assert.Equal(t, []prtrie.Result{
		{Term: "tester", Freq: 10},
		{Term: "testing", Freq: 5},
	}, p.TopKForPrefix("", 2), "expected all top terms")
	assert.Equal(t, []prtrie.Result{
		{Term: "tester", Freq: 10},
		{Term: "testing", Freq: 5},
	}, p.TopKForPrefix("test", 2), "expected all top terms")
	assert.Equal(t, []prtrie.Result{
		{Term: "tester", Freq: 10},
	}, p.TopKForPrefix("teste", 2), "expected all top terms")
	assert.Equal(t, []prtrie.Result{
		{Term: "testing", Freq: 5},
	}, p.TopKForPrefix("testin", 2), "expected all top terms")

	p.AddTerm("terrarium", 7)

	assert.Equal(t, []prtrie.Result{
		{Term: "tester", Freq: 10},
		{Term: "terrarium", Freq: 7},
	}, p.TopKForPrefix("", 2), "expected all top terms")
	assert.Equal(t, []prtrie.Result{
		{Term: "tester", Freq: 10},
		{Term: "terrarium", Freq: 7},
	}, p.TopKForPrefix("te", 2), "expected all top terms")
	assert.Equal(t, []prtrie.Result{
		{Term: "tester", Freq: 10},
		{Term: "testing", Freq: 5},
	}, p.TopKForPrefix("tes", 2), "expected all top terms")
	assert.Equal(t, []prtrie.Result{
		{Term: "terrarium", Freq: 7},
	}, p.TopKForPrefix("ter", 2), "expected all top terms")
}

func TestPrtrieAddingOverlappingTermsThenAddCommon(t *testing.T) {
	p := prtrie.NewPruningRadixTrie()
	p.AddTerm("testing", 5)
	p.AddTerm("tester", 10)
	assert.Equal(t, uint64(2), p.GetTotalTermCount())
	p.AddTerm("test", 80)
	assert.Equal(t, uint64(3), p.GetTotalTermCount())
	assert.Equal(t, []prtrie.Result{
		{Term: "test", Freq: 80},
		{Term: "tester", Freq: 10},
		{Term: "testing", Freq: 5},
	}, p.TopKForPrefix("test", 4), "expected all top terms")
}

func TestPrtrieGetManyChildrenTermsSameFrequency(t *testing.T) {
	p := prtrie.NewPruningRadixTrie()
	p.AddTerm("abc", 1)
	p.AddTerm("abd", 1)
	p.AddTerm("abe", 1)
	p.AddTerm("abcd", 1)
	p.AddTerm("abce", 1)
	p.AddTerm("abcde", 1)
	p.AddTerm("abcdef", 1)
	p.AddTerm("abcdefg", 1)
	p.AddTerm("abcdefh", 1)
	// TODO: flaky test, if we change the ResultSet results with equal frequency might
	// be sorted differently
	assert.Equal(t, []prtrie.Result{
		// {Term: "abcd", Freq: 1},
		// {Term: "abc", Freq: 1},
		{Term: "abc", Freq: 1},
		{Term: "abcd", Freq: 1},
	}, p.TopKForPrefix("ab", 2), "expected all top terms")
}

func TestPrtrieString(t *testing.T) {
	p := prtrie.NewPruningRadixTrie()
	assert.Equal(t, "[0,0]\n", p.String())
	p.AddTerm("test", 1)
	assert.Equal(t, strings.Join([]string{
		"[0,1]",
		"test[1,1]",
		"",
	}, "\n"), p.String())
	p.AddTerm("testing", 2)
	assert.Equal(t, strings.Join([]string{
		"[0,2]",
		"test[1,2]",
		"    ing[2,2]",
		"",
	}, "\n"), p.String())
	p.AddTerm("tester", 5)
	assert.Equal(t, strings.Join([]string{
		"[0,5]",
		"test[1,5]",
		"    er[5,5]",
		"    ing[2,2]",
		"",
	}, "\n"), p.String())
	assert.Equal(t, []prtrie.Result{
		{Term: "tester", Freq: 5},
		{Term: "testing", Freq: 2},
	}, p.TopKForPrefix("tes", 2), "expected all top terms")
	p.AddTerm("food", 1)
	assert.Equal(t, strings.Join([]string{
		"[0,5]",
		"test[1,5]",
		"    er[5,5]",
		"    ing[2,2]",
		"food[1,1]",
		"",
	}, "\n"), p.String())
	p.AddTerm("food", 12)
	assert.Equal(t, []prtrie.Result{
		{Term: "food", Freq: 13},
		{Term: "tester", Freq: 5},
		{Term: "testing", Freq: 2},
		{Term: "test", Freq: 1},
	}, p.TopKForPrefix("", 4), "expected all top terms")
	assert.Equal(t, strings.Join([]string{
		"[0,13]",
		"food[13,13]",
		"test[1,5]",
		"    er[5,5]",
		"    ing[2,2]",
		"",
	}, "\n"), p.String())
	p.AddTerm("testing", 48)
	assert.Equal(t, []prtrie.Result{
		{Term: "testing", Freq: 50},
		{Term: "tester", Freq: 5},
	}, p.TopKForPrefix("tes", 2), "expected all top terms")
	assert.Equal(t, strings.Join([]string{
		"[0,50]",
		"test[1,50]",
		"    ing[50,50]",
		"    er[5,5]",
		"food[13,13]",
		"",
	}, "\n"), p.String())
}

func TestResultString(t *testing.T) {
	assert.Equal(t, "test:5", prtrie.Result{Term: "test", Freq: 5}.String())
}
