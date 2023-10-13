package pruningradixtrie

import (
	"fmt"
	"sort"
	"strings"
)

type Result struct {
	Term string
	Freq uint64
}

type ResultSet interface {
	Results() []Result
}

type PruningRadixTrie interface {
	AddTerm(term string, count uint64)
	TopKForPrefix(prefix string, k int) []Result
	// FindAllTForPrefix(prefix string) []Result
}

type pruningRadixTrie struct {
	trie      *node
	termCount uint64
}

var _ PruningRadixTrie = &pruningRadixTrie{}
var _ fmt.Stringer = &pruningRadixTrie{}

func NewPruningRadixTrie() *pruningRadixTrie {
	return &pruningRadixTrie{
		trie: newNode("", 0),
	}
}

func (p *pruningRadixTrie) GetTotalTermCount() uint64 {
	return p.termCount
}

func (p *pruningRadixTrie) AddTerm(term string, count uint64) {
	if term == "" || count == 0 {
		return
	}
	var list []*node
	p.addTerm(p.trie, term, count, 0, list)
}

func (p *pruningRadixTrie) addTerm(
	cur *node,
	term string,
	count uint64,
	level int,
	list []*node,
) {
	list = append(list, cur)

	for i, child := range cur.children {
		key := child.key
		common := findCommon(key, term)
		if common > 0 {
			if common == len(term) && common == len(key) {
				//term already existed
				//existing ab
				//new      ab
				if child.count == 0 {
					p.termCount++
				}
				child.count += count
				newMax := max(child.count, child.maxChildCount)
				child.maxChildCount = newMax
				// todo: see if we can only update nodes in path
				// where max count was updated
				// list = append(list, child)
				// updateMaxCounts(list, child.count)
				reorder(list, newMax)

			} else if common == len(term) {
				//new is subkey
				//existing abcd
				//new      ab
				//if new is shorter (== common), then node(count) and only 1. children add (clause2)
				// 0123 -> [:2]:[2:]
				// abcd -> ab:cd

				// split with unicode
				// https://stackoverflow.com/a/56129336
				prefix := term
				suffix := key[common:]
				oldChildCount := child.count
				newChild := newNode(prefix, count)
				newMax := max(cur.count, cur.maxChildCount)
				newChild.maxChildCount = newMax
				cur.children[i] = newChild
				// updateChildren(cur)
				reorder(list, newMax)
				// adding the old child back under the split term as if it was a new term added
				p.addTerm(newChild, suffix, oldChildCount, level+1, list)
			} else if common == len(key) {
				//if oldkey shorter (==common), then recursive addTerm (clause1)
				//existing: te
				//new:      test
				p.addTerm(child, term[common:], count, level+1, list)
			} else {
				//old and new have common substrings
				//existing: test
				//new:      team
				commonPrefix := term[:common]
				termSuffix := term[common:]
				splitSuffix := key[common:]
				// create a new node of the common substrings
				// the count is zero since the intersection was never explicitly added prior
				split := newNode(commonPrefix, 0)
				split.children = append(split.children, child)
				split.maxChildCount = max(child.maxChildCount, child.count, count)
				child.key = splitSuffix

				cur.children[i] = split
				p.addTerm(split, termSuffix, count, level+1, list)
			}
			return
		}
	}

	newChild := newNode(term, count)
	cur.children = append(cur.children, newChild)
	list = append(list, newChild)
	updateMaxCounts(list, count)
	updateChildren(cur)
	p.termCount++
}

func findCommon(key, term string) int {
	common := 0
	// TODO: make it unicode compatible.
	for i := 0; i < min(len(key), len(term)); i++ {
		if key[i] != term[i] {
			break
		}
		common++
	}
	return common
}

func reorder(nodes []*node, count uint64) {
	for i := len(nodes) - 1; i >= 0; i-- {
		n := nodes[i]
		if count > n.maxChildCount {
			n.maxChildCount = count
		}
		updateChildren(n)
	}
}

// updateChildren sort children of a node in descending order
func updateChildren(n *node) {
	children := n.children
	sort.Slice(children, func(i, j int) bool {
		return children[j].maxChildCount < children[i].maxChildCount
	})
}

func updateMaxCounts(nodes []*node, count uint64) {
	// TODO: check if update max left to root is more correct
	for _, n := range nodes {
		if count > n.maxChildCount {
			n.maxChildCount = count
		}
	}
}

// TopKForPrefix implements PruningRadixTrie.
func (p *pruningRadixTrie) TopKForPrefix(prefix string, k int) []Result {
	if k <= 0 {
		return nil
	}
	var results []Result
	return topKForPrefix(prefix, "", p.trie, k, results)
}

// TopKForPrefix implements PruningRadixTrie.
func topKForPrefix(prefix, path string, cur *node, k int, results []Result) []Result {
	if len(results) == k && cur.maxChildCount <= results[k-1].Freq {
		return results
	}

	noPrefix := prefix == ""
	if len(cur.children) == 0 {
		return results
	}
	for _, child := range cur.children {
		key := child.key
		if len(results) == k && child.count <= results[k-1].Freq && child.maxChildCount <= results[k-1].Freq {
			if noPrefix {
				continue
			}
			break
		}
		if noPrefix || strings.HasPrefix(key, prefix) {
			if child.count > 0 {
				results = insertTopKSuggestion(path+key, child.count, k, results)
			}
			if len(child.children) > 0 {
				results = topKForPrefix("", path+key, child, k, results)
			}
			if !noPrefix {
				break
			}
		} else if strings.HasPrefix(prefix, key) {
			if len(child.children) > 0 {
				return topKForPrefix(prefix[len(key):], path+key, child, k, results)
			}
		}
	}
	return results
}

// TODO: switch to heap for insertion, and heapsort for results
func insertTopKSuggestion(term string, freq uint64, k int, results []Result) []Result {
	//at the end/highest index is the lowest value
	// >  : old take precedence for equal rank
	// >= : new take precedence for equal rank
	if len(results) > k || (len(results) == k && freq < results[k-1].Freq) {
		return results
	}
	result := Result{Term: term, Freq: freq}
	results = append(results, result)
	sort.Slice(results, func(i, j int) bool {
		return results[j].Freq < results[i].Freq
	})

	if len(results) > k {
		results = results[:k]
	}
	return results
}

// String the string representation of the trie in tree like format
// [0, 777]
//
// qux[777, 777]
// bar[77, 77]
// foo[7, 7]
func (p *pruningRadixTrie) String() string {
	s := []struct {
		node  *node
		level int
	}{{node: p.trie, level: 0}}
	b := strings.Builder{}
	for len(s) > 0 {
		top := s[len(s)-1]
		n := top.node
		level := top.level
		s = s[:len(s)-1]
		b.WriteString(fmt.Sprintf(
			"%s%s[%d,%d]\n",
			strings.Repeat(" ", level),
			n.key,
			n.count,
			n.maxChildCount,
		))
		for i := len(n.children) - 1; i >= 0; i-- {
			s = append(s, struct {
				node  *node
				level int
			}{node: n.children[i], level: len(n.key)})
		}
	}
	return b.String()
}

func (r Result) String() string {
	return fmt.Sprintf("%s:%d", r.Term, r.Freq)
}
