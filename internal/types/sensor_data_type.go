package types

import "time"

type SensorDataType struct {
	ID        int
	CreatedAt time.Time
	UpdatedAt time.Time
	Name      string
}
