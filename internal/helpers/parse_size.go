package helpers

import (
	"errors"
	"math"
	"strconv"
	"strings"
)

func hasFraction(x float64) bool {
	_, frac := math.Modf(x)
	return frac != 0
}

func ParseSize(s string) (int, error) {
	s = strings.TrimSpace(s)
	s = strings.ToLower(s)

	if s == "" || s[0] < '0' || s[0] > '9' {
		return 0, errors.New("invalid size format")
	}

	var size float64
	var unit string

	i := 0
	for ; i < len(s); i++ {
		if s[i] >= 'a' && s[i] <= 'z' {
			break
		}
	}

	size, _ = strconv.ParseFloat(s[:i], 64)
	unit = s[i:]

	switch unit {
	case "b", "":
		if hasFraction(size) {
			return 0, errors.New("wrong size")
		}
		return int(size), nil
	case "kb":
		return int(size * (1 << 10)), nil
	case "mb":
		return int(size * (1 << 20)), nil
	case "gb":
		return int(size * (1 << 30)), nil
	default:
		return 0, errors.New("invalid unit")
	}

}
