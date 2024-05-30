package utils

func IndexOf(value any, options []any) int {
	for i, option := range options {
		if option == value {
			return i
		}
	}
	return 0
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
