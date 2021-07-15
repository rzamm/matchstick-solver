// Package ec (Enumerative Combinatorics) extends the functionality of gonum.
package ec

// NextCombination generates the combination after s, overwriting the slice.
// Ported from "gonum.org/v1/gonum/stat/combin".
func NextCombination(s []int, n, k int) {
	for j := k - 1; j >= 0; j-- {
		if s[j] == n+j-k {
			continue
		}
		s[j]++
		for l := j + 1; l < k; l++ {
			s[l] = s[j] + l - j
		}
		break
	}
}
