package parkings

import (
	"errors"

	"gorm.io/gorm"
)

var ErrZoneNotFound = errors.New("parking zone not found")

type Repository interface {
	Create(zone *Zone) error
	GetAll() ([]*ZoneWithSpots, error)
	GetByID(zoneId uint) (*ZoneWithSpots, error)
}

// Data wrapper layout to hold the calculated dynamic metrics cleanly
type ZoneWithSpots struct {
	Zone
	AvailableSpots int `gorm:"column:available_spots"`
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(zone *Zone) error {
	return r.db.Create(zone).Error
}

func (r *repository) GetAll() ([]*ZoneWithSpots, error) {
	var zones []*Zone
	// Fetch zones standardly without looking at reservations
	err := r.db.Find(&zones).Error
	if err != nil {
		return nil, err
	}

	var result []*ZoneWithSpots
	for _, z := range zones {
		result = append(result, &ZoneWithSpots{
			Zone:           *z,
			AvailableSpots: z.TotalCapacity, // Temporary fallback: all spots available
		})
	}
	return result, nil
}

func (r *repository) GetByID(zoneId uint) (*ZoneWithSpots, error) {
	var zone Zone
	err := r.db.First(&zone, zoneId).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrZoneNotFound
		}
		return nil, err
	}

	return &ZoneWithSpots{
		Zone:           zone,
		AvailableSpots: zone.TotalCapacity, // Temporary fallback
	}, nil
}

// func (r *repository) GetAll() ([]*ZoneWithSpots, error) {
// 	var zones []*ZoneWithSpots

// 	// Subquery to count active reservations for each zone
// 	subQuery := r.db.Table("reservations").
// 		Select("COUNT(*)").
// 		Where("reservations.zone_id = zones.id AND reservations.status = 'active'")

// 	err := r.db.Table("zones").
// 		Select("zones.*, (zones.total_capacity - (?)) as available_spots", subQuery).
// 		Where("zones.deleted_at IS NULL").
// 		Find(&zones).Error

// 	if err != nil {
// 		return nil, err
// 	}
// 	return zones, nil
// }

// func (r *repository) GetByID(zoneId uint) (*ZoneWithSpots, error) {
// 	var zone ZoneWithSpots

// 	subQuery := r.db.Table("reservations").
// 		Select("COUNT(*)").
// 		Where("reservations.zone_id = zones.id AND reservations.status = 'active'")

// 	err := r.db.Table("zones").
// 		Select("zones.*, (zones.total_capacity - (?)) as available_spots", subQuery).
// 		Where("zones.id = ? AND zones.deleted_at IS NULL", zoneId).
// 		First(&zone).Error

// 	if err != nil {
// 		if errors.Is(err, gorm.ErrRecordNotFound) {
// 			return nil, ErrZoneNotFound
// 		}
// 		return nil, err
// 	}
// 	return &zone, nil
// }
