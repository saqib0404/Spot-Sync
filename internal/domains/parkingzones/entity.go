package parkings

import (
	"Spot-Sync/internal/domains/parkingzones/dto"
	"time"

	"gorm.io/gorm"
)

type ZoneType string

const (
	ZoneGeneral ZoneType = "general"
	ZoneEV      ZoneType = "ev_charging"
	ZoneCovered ZoneType = "covered"
)

type Zone struct {
	gorm.Model
	Name          string   `json:"name" gorm:"type:varchar(255);not null"`
	Type          ZoneType `json:"type" gorm:"type:varchar(50);not null;default:'general'"`
	TotalCapacity int      `json:"total_capacity" gorm:"not null"`
	PricePerHour  float64  `json:"price_per_hour" gorm:"type:decimal(10,2);not null"`
}

func (z *Zone) ToResponse() *dto.ZoneResponse {
	return &dto.ZoneResponse{
		ID:            z.ID,
		Name:          z.Name,
		Type:          string(z.Type),
		TotalCapacity: z.TotalCapacity,
		PricePerHour:  z.PricePerHour,
		CreatedAt:     z.CreatedAt.Format(time.RFC3339),
		UpdatedAt:     z.UpdatedAt.Format(time.RFC3339),
	}
}
