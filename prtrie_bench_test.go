package pruningradixtrie_test

import (
	"archive/zip"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
	"testing"

	prtrie "github.com/elielamora/pruningradixtrie"
)

func BenchmarkTrivialOperation(b *testing.B) {
	t := 0
	for i := 0; i < b.N; i++ {
		t += 1
	}
}

func BenchmarkPruningRadixTrie_Insert(b *testing.B) {
	entries, err := getEntries()
	if err != nil {
		b.Fatalf("unexpected error: %s", err)
	}
	p := prtrie.NewPruningRadixTrie()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		entry := entries[i%len(entries)]
		p.AddTerm(entry.Term, entry.Freq)
	}
}

func BenchmarkPruningRadixTrie_TopKMicrosoft(b *testing.B) {
	entries, err := getEntries()
	if err != nil {
		b.Fatalf("unexpected error: %s", err)
	}
	p := prtrie.NewPruningRadixTrie()
	for _, entry := range entries {
		p.AddTerm(entry.Term, entry.Freq)
	}
	b.ResetTimer()
	searchString := "microsoft"
	for _, k := range []int{10, 50, 100, 1000, 10000} {
		for i := 1; i <= len(searchString); i++ {
			prefix := searchString[:i]
			b.Run(fmt.Sprintf("Top%d: %s", k, prefix), func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					p.TopKForPrefix(prefix, k)
				}
			})
		}
	}
}

func _TestVisualize(t *testing.T) {
	entries, err := getEntries()
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	p := prtrie.NewPruningRadixTrie()
	for _, entry := range entries {
		p.AddTerm(entry.Term, entry.Freq)
	}
	fmt.Println(p.String())
}

func getEntries() ([]prtrie.Result, error) {
	termsString, err := unzip()
	if err != nil {
		return nil, err
	}
	terms := strings.Split(termsString, "\r\n")
	var results []prtrie.Result
	for _, term := range terms {
		if term == "" {
			continue
		}
		x := strings.Split(term, "\t")
		if len(x) != 2 {
			return nil, errors.New(fmt.Sprintf("unexpected format, expected 2 tab separated values but got %d: %s", len(x), term))
		}
		freq, err := strconv.ParseUint(x[1], 10, 64)
		if err != nil {
			return nil, err
		}
		results = append(results, prtrie.Result{
			Term: x[0],
			Freq: freq,
		})
	}
	return results, nil
}

// unzip unzips the bundled test terms file
func unzip() (string, error) {
	r, err := zip.OpenReader("./terms.zip")
	if err != nil {
		return "", errors.Join(errors.New("tried opening the zip reader"), err)
	}

	for _, f := range r.File {
		if f.Name != "terms.txt" {
			return "", errors.New(fmt.Sprintf(
				"expected 'terms.txt' but got %s",
				f.Name,
			))
		}
		rc, err := f.Open()
		if err != nil {
			return "", err
		}
		defer rc.Close()

		strbuilder := &strings.Builder{}
		_, err = io.Copy(strbuilder, rc)
		if err != nil {
			return "", err
		}
		return strbuilder.String(), nil
	}
	return "", nil
}
