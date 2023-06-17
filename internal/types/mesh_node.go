package types

import (
	"time"
)

type MeshNode struct {
	UUID      UUID              `json:"uuid"`
	UpdateID  *MeshNodeUpdateID `json:"updateId,omitempty"`
	CreatedAt time.Time         `json:"createdAt"`
	UpdatedAt *time.Time        `json:"updatedAt,omitempty"`
	Latitude  float32           `json:"latitude"`
	Longitude float32           `json:"longitude"`
}
