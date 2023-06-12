package types

import (
	"time"

	"github.com/mdma-backend/mdma-backend/internal/types/permission"
)

type RoleID uint

type Role struct {
	ID          RoleID                  `json:"id,omitempty"`
	CreatedAt   time.Time               `json:"createAt"`
	UpdatedAt   *time.Time              `json:"updatedAt,omitempty"`
	Name        string                  `json:"name"`
	Permissions []permission.Permission `json:"permissions,omitempty"`
}
