package utils

func MinMaxIndex(index int, max int) int {
	if index < 0 {
		return max - 1
	}
	if index >= max {
		return 0
	}
	return index
}
