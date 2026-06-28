package reservations

import (
	"Spot-Sync/internal/domains/parkings"

	"gorm.io/gorm"
)

type ReservationStatus string

const (
	StatusActive    ReservationStatus = "active"
	StatusCompleted ReservationStatus = "completed"
	StatusCancelled ReservationStatus = "cancelled"
)

type Reservation struct {
	gorm.Model
	UserID       uint              `json:"user_id" gorm:"not null"`
	ZoneID       uint              `json:"zone_id" gorm:"not null"`
	Zone         parkings.Zone     `json:"zone" gorm:"foreignKey:ZoneID"`
	LicensePlate string            `json:"license_plate" gorm:"type:varchar(15);not null"`
	Status       ReservationStatus `json:"status" gorm:"type:varchar(50);not null;default:'active'"`
}
