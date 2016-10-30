package router

import (
	"strings"
)

// min return minimum of two integers
func min(a, b int) int {
	if a <= b {
		return a
	}

	return b
}

// lcp return length of longest common prefix of two strings
func lcp(s1, s2 string) int {
	len := min(len(s1), len(s2))
	for i := 0; i < len; i++ {
		if s1[i] != s2[i] {
			return i
		}
	}

	return len
}

// findFirstParam returns the index of the first instance of symbol in path and it's length,
// or -1, 0 if symbol is not present in path.
func findFirstParam(path string, symbol byte) (pos int, length int) {
	pos = strings.IndexByte(path, symbol)

	if pos == -1 {
		return -1, 0
	}

	pathLen := len(path)

	id := pos
	for id < pathLen && path[id] != '/' {
		id++
	}

	return pos, id - pos
}
