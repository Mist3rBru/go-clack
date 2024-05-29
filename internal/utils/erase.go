package utils

func EraseCurrentLine() string {
	return "\x1b[K"
}

func EraseDown() string {
	return "\x1b[J"
}