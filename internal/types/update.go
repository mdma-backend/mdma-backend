package types

import "time"

type Update struct {
	ID        int
	CreatedAt time.Time
	UpdatedAt time.Time
	Version   string
	Data      []byte
}
