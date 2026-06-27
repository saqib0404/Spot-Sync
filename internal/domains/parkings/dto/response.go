package dto

type ZoneResponse struct {
	ID             uint    `json:"id"`
	Name           string  `json:"name"`
	Type           string  `json:"type"`
	TotalCapacity  int     `json:"total_capacity"`
	AvailableSpots *int    `json:"available_spots,omitempty"`
	PricePerHour   float64 `json:"price_per_hour"`
	CreatedAt      string  `json:"created_at"`
	UpdatedAt      string  `json:"updated_at,omitempty"`
}

// Global API envelope matching response schema
type APIEnvelope struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}
