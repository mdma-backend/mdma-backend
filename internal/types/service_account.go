package types

import "time"

type ServiceAccount struct {
	ID        int
	Role      Role
	CreatedAt time.Time
	UpdatedAt time.Time
	Name      string
}
