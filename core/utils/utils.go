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
	ansiRegex := regexp.MustCompile(`[\\u001B\\u009B][[\\]()#;?]*(?:(?:(?:[a-zA-Z\\d]*(?:;[-a-zA-Z\\d\\/#&.:=?%@~_]*)*)?\\u0007)|(?:(?:\\d{1,4}(?:;\\d{0,4})*)?[\\dA-PR-TZcf-ntqry=><~]))`)
	parsedStr := ansiRegex.ReplaceAllString(str, "")
	length := 0

	for i := 0; i < len(parsedStr); i++ {
		r := rune(parsedStr[i])
		if isControlCharacter(r) || isCombiningCharacter(r) {
			continue
		} else if isSurrogatePair(r) {
			i++
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
