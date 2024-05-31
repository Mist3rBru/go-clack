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

func Max(n1, n2 int) int {
	if n1 > n2 {
		return n1
	}
	return n2
}

func Min(n1, n2 int) int {
	if n1 < n2 {
		return n1
	}
	return n2
}
