package util

import "strconv"

func StringToInt(text string, defaultValue int) int {
	value, err := strconv.Atoi(text)

	if err != nil {
		return defaultValue
	}

	return value
}
