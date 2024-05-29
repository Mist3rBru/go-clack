package utils

import (
	"math"
	"strings"
)

func DiffLines(a string, b string) []int {
	diff := []int{}

	if a == b {
		return diff
	}

	aLines := strings.Split(a, "\n")
	bLines := strings.Split(b, "\n")
	for i := range int(math.Max(float64(len(aLines)), float64(len(bLines)))) {
		if i+1 > len(aLines) || i+1 > len(bLines) || aLines[i] != bLines[i] {
			diff = append(diff, i)
		}
	}

	return diff
}
