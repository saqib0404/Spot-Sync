package dto

type CreateReservationRequest struct {
	ZoneID       uint   `json:"zone_id" validate:"required,gt=0"`
	LicensePlate string `json:"license_plate" validate:"required,min=2,max=15"` // Constrained to 15 max characters
}

type NestedZoneInfo struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}
