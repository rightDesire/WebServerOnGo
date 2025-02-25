package utils

import (
	"strconv"
)

func StringToInt(str string) (uint, error) {
	num, err := strconv.Atoi(str)
	if err != nil {
		return 0, err
	}
	return uint(num), nil
}
