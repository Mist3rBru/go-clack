package utils

import (
	"strings"
)

func DiffLines(a string, b string) []int {
	diff := []int{}

	if a == b {
		return diff
	}

	aLines := strings.Split(a, "\n")
	bLines := strings.Split(b, "\n")
	for i := range Max(len(aLines), len(bLines)) {
		if i >= len(aLines) || i >= len(bLines) || aLines[i] != bLines[i] {
			diff = append(diff, i)
		}
	}

	return diff
}
