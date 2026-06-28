package dto

type ReservationResponse struct {
	ID           uint            `json:"id"`
	UserID       uint            `json:"user_id,omitempty"`
	ZoneID       uint            `json:"zone_id,omitempty"`
	LicensePlate string          `json:"license_plate"`
	Status       string          `json:"status"`
	Zone         *NestedZoneInfo `json:"zone,omitempty"`
	CreatedAt    string          `json:"created_at"`
	UpdatedAt    string          `json:"updated_at,omitempty"`
}

type ReservationAPIEnvelope struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}
