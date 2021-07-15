package ec

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"gonum.org/v1/gonum/stat/combin"
)

func BenchmarkCombin(b *testing.B) {
	n := 30
	k := 10
	s1 := make([]int, k)
	s2 := make([]int, k)

	library := func() {
		gen := combin.NewCombinationGenerator(n, k)
		idx := 0
		for gen.Next() {
			gen.Combination(s1)
			idx++
		}
	}

	manual := func() {
		// init
		for i := 0; i < k; i++ {
			s2[i] = i
		}
		idx := combin.CombinationIndex(s2, n, k)
		total := combin.Binomial(n, k)

		// do first one
		idx++

		for idx < total {
			NextCombination(s2, n, k)
			idx++
		}
	}

	// run both functions, compare which is faster
	for i := 0; i < b.N; i++ {
		library()
		manual()
	}
}

func TestCombin(t *testing.T) {
	n := 25
	k := 4
	s := make([]int, k)
	gen := combin.NewCombinationGenerator(n, k)
	idx := 0
	for gen.Next() {
		fmt.Println(idx, gen.Combination(s))
		idx++
	}
	assert.Equal(t, 12650, idx)
	assert.ElementsMatch(t, []int{21, 22, 23, 24}, s)
}

func TestCombinManual(t *testing.T) {
	n := 25
	k := 4
	s := make([]int, k)

	// init s
	for i := 0; i < k; i++ {
		s[i] = i
	}
	s = []int{19, 20, 21, 22}

	idx := combin.CombinationIndex(s, n, k)
	total := combin.Binomial(n, k)

	for idx < total {
		fmt.Println(idx, s)
		NextCombination(s, n, k)
		idx++
	}
	assert.Equal(t, 12650, idx)
	assert.ElementsMatch(t, []int{21, 22, 23, 24}, s)
}
