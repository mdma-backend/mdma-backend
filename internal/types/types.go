package types

import (
	"errors"
	"fmt"
	"strconv"
)

type ID interface {
	~uint
}

func IDFromString[T ID](s string) (T, error) {
	id, err := strconv.Atoi(s)
	if err != nil {
		return 0, fmt.Errorf("%s is not an integer", s)
	}

	if id <= 0 {
		return 0, errors.New("id must be greater than 0")
	}

	return T(uint(id)), nil
}
