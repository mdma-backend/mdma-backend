package types

import "time"

type SensorData struct {
	ID         [8]byte
	Controller Controller
	Type       SensorDataType
	CreatedAt  time.Time
	MeasuredAt time.Time
	Value      string
}
