package types

import "time"

type ServiceAccountID uint

type ServiceAccount struct {
	ID        ServiceAccountID `json:"id,omitempty"`
	RoleID    RoleID           `json:"roleId,omitempty"`
	CreatedAt time.Time        `json:"createdAt"`
	UpdatedAt *time.Time       `json:"updatedAt,omitempty"`
	Username  string           `json:"username"`
	Token     string           `json:"token"`
}
