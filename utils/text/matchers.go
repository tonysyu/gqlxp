package text

import "strings"

// SuggestSimilarWords returns words from allowedWords with edit distance <= 2 from word.
func SuggestSimilarWords(word string, allowedWords []string) []string {
	lower := strings.ToLower(word)
	var suggestions []string
	for _, allowed := range allowedWords {
		if levenshtein(lower, allowed) <= 2 {
			suggestions = append(suggestions, allowed)
		}
	}
	return suggestions
}

// levenshtein returns the edit distance between ASCII strings a and b.
func levenshtein(a, b string) int {
	m, n := len(a), len(b)
	if m == 0 {
		return n
	}
	if n == 0 {
		return m
	}
	prev := make([]int, n+1)
	curr := make([]int, n+1)
	for j := range prev {
		prev[j] = j
	}
	for i := range m {
		curr[0] = i + 1
		for j := range n {
			cost := 1
			if a[i] == b[j] {
				cost = 0
			}
			curr[j+1] = min(curr[j]+1, min(prev[j+1]+1, prev[j]+cost))
		}
		prev, curr = curr, prev
	}
	return prev[n]
}
