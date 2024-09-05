package utils

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

	length := 0
	inEscapeCode := false

	for i := 0; i < len(str); i++ {
		r := rune(str[i])

		if inEscapeCode {
			if r == 'm' {
				inEscapeCode = false
			}
			continue
		}

		if r == '\x1b' {
			inEscapeCode = true
			// length++ // count the escape code as 1 character
			continue
		}

		if isControlCharacter(r) || isCombiningCharacter(r) {
			continue
		}

		if isSurrogatePair(r) {
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
