package types

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/gofrs/uuid"
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

type UUID struct {
	uuid.UUID
}

func (id UUID) MarshalJSON() ([]byte, error) {
	return json.Marshal(id.String())
}

func (id *UUID) UnmarshalJSON(b []byte) error {
	var v interface{}
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}

	switch value := v.(type) {
	case string:
		idx, err := uuid.FromString(value)
		if err != nil {
			return err
		}
		id.UUID = idx
		return nil
	default:
		return errors.New("invalid uuid format")
	}
}

func UUIDFromString(s string) (UUID, error) {
	id, err := uuid.FromString(s)
	if err != nil {
		return UUID{uuid.Nil}, err
	}

	return UUID{id}, nil
}
