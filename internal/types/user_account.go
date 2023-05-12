package types

import "time"

type UserAccount struct {
	ID        int
	Role      Role
	CreatedAt time.Time
	UpdatedAt time.Time
	Username  string
	Password  []byte
}
