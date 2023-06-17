package types

import (
	"time"
)

type MeshNodeUpdateID uint

type MeshNodeUpdate struct {
	ID        RoleID    `json:"id,omitempty"`
	CreatedAt time.Time `json:"createAt"`
	Version   string    `json:"version"`
	Data      []byte    `json:"data,omitempty"`
}
