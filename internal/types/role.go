package types

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/mdma-backend/mdma-backend/internal/types/permission"
)

type RoleID int

func RoleIDFromString(s string) (RoleID, error) {
	id, err := strconv.Atoi(s)
	if err != nil {
		return 0, fmt.Errorf("%s is not an integer", s)
	}

	if id <= 0 {
		return 0, errors.New("id must be greater than 0")
	}

	return RoleID(id), nil
}

type Role struct {
	ID          RoleID                  `json:"id,omitempty"`
	CreatedAt   time.Time               `json:"createAt"`
	UpdatedAt   *time.Time              `json:"updatedAt,omitempty"`
	Name        string                  `json:"name"`
	Permissions []permission.Permission `json:"permissions,omitempty"`
}
