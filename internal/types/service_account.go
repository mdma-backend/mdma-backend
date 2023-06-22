package types

import "time"

type ServiceAccountID uint

type ServiceAccount struct {
	ID        ServiceAccountID `json:"id,omitempty"`
	RoleID    *RoleID          `json:"roleId,omitempty"`
	CreatedAt time.Time        `json:"createdAt"`
	UpdatedAt *time.Time       `json:"updatedAt,omitempty"`
	Name      string           `json:"name"`
	Token     string           `json:"token,omitempty"`
}
