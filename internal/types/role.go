package types

import "time"

type Role struct {
	ID          int
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Name        string
	Permissions []string
}
