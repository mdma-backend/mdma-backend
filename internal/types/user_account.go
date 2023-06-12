package types

import "time"

type UserAccountID uint

type UserAccount struct {
	ID        UserAccountID `json:"id,omitempty"`
	RoleID    *RoleID       `json:"roleId,omitempty"`
	CreatedAt time.Time     `json:"createdAt"`
	UpdatedAt *time.Time    `json:"updatedAt,omitempty"`
	Username  string        `json:"username"`
	Password  string        `json:"password,omitempty"`
}
