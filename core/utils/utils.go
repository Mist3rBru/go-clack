package utils

import (
	"regexp"
	"strings"
)

func DiffLines(a string, b string) []int {
	diff := []int{}

	if a == b {
		return diff
	}

	aLines := strings.Split(a, "\n")
	bLines := strings.Split(b, "\n")
	for i := range max(len(aLines), len(bLines)) {
		if i >= len(aLines) || i >= len(bLines) || aLines[i] != bLines[i] {
			diff = append(diff, i)
		}
	}

	return diff
}

func isControlCharacter(r rune) bool {
	return r <= 0x1f || (r >= 0x7f && r <= 0x9f)
}

func isCombiningCharacter(r rune) bool {
	return r >= 0x300 && r <= 0x36f
}

func isSurrogatePair(r rune) bool {
	return r >= 0xd800 && r <= 0xdbff
}

func StrLength(str string) int {
	if len(str) == 0 {
		return 0
	}
	ansiRegex := regexp.MustCompile(`\x1b\[[0-9;]*m`)
	parsedStr := ansiRegex.ReplaceAllString(str, "")
	length := 0

	for _, r := range parsedStr {
		if isControlCharacter(r) || isCombiningCharacter(r) || isSurrogatePair(r) {
			continue
		}
		length++
	}

	return length
}

func MinMaxIndex(index int, max int) int {
	if index < 0 {
		return max - 1
	}
	if index >= max {
		return 0
	}
	return index
}
