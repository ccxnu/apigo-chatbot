package domain

import "time"

// Base contains common fields required in all API requests
// This ensures traceability and audit trail for all operations
type Base struct {
	IdSession     string    `json:"idSession" validate:"required" doc:"Session identifier"`
	IdRequest     string    `json:"idRequest" validate:"required,uuid4" doc:"Unique request ID (UUID v4)"`
	Process       string    `json:"process" validate:"required" doc:"Process name initiating the request"`
	IdDevice      string    `json:"idDevice" validate:"required" doc:"Device identifier"`
	DeviceAddress string    `json:"deviceAddress" validate:"required,ip" doc:"Device IP address"`
	DateProcess   time.Time `json:"dateProcess" validate:"required" doc:"Timestamp when request was initiated"`
}
