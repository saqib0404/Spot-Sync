package reservations

import (
	"errors"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var (
	ErrZoneFull            = errors.New("parking zone is completely full")
	ErrReservationNotFound = errors.New("reservation not found")
)

type Repository interface {
	CreateAtomic(userID uint, zoneID uint, licensePlate string) (*Reservation, error)
	GetByUserID(userID uint) ([]*Reservation, error)
	GetByID(id uint) (*Reservation, error)
	UpdateStatus(reservation *Reservation, status ReservationStatus) error
	GetAllAdmin() ([]*Reservation, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) CreateAtomic(userID uint, zoneID uint, licensePlate string) (*Reservation, error) {
	var reservation *Reservation

	err := r.db.Transaction(func(tx *gorm.DB) error {
		var zone struct {
			ID            uint
			TotalCapacity int
		}

		// Lock row to protect concurrent reads
		if err := tx.Table("zones").
			Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("id = ? AND deleted_at IS NULL", zoneID).
			First(&zone).Error; err != nil {
			return err
		}

		// Count only currently active reservations
		var activeCount int64
		if err := tx.Model(&Reservation{}).
			Where("zone_id = ? AND status = 'active'", zoneID).
			Count(&activeCount).Error; err != nil {
			return err
		}

		if int(activeCount) >= zone.TotalCapacity {
			return ErrZoneFull
		}

		reservation = &Reservation{
			UserID:       userID,
			ZoneID:       zoneID,
			LicensePlate: licensePlate,
			Status:       StatusActive,
		}

		return tx.Create(reservation).Error
	})

	if err != nil {
		return nil, err
	}
	return reservation, nil
}

func (r *repository) GetByUserID(userID uint) ([]*Reservation, error) {
	var list []*Reservation
	err := r.db.Preload("Zone").Where("user_id = ? AND deleted_at IS NULL", userID).Find(&list).Error
	return list, err
}

func (r *repository) GetByID(id uint) (*Reservation, error) {
	var item Reservation
	err := r.db.Where("id = ? AND deleted_at IS NULL", id).First(&item).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrReservationNotFound
		}
		return nil, err
	}
	return &item, nil
}

func (r *repository) UpdateStatus(res *Reservation, status ReservationStatus) error {
	return r.db.Model(res).Update("status", status).Error
}

func (r *repository) GetAllAdmin() ([]*Reservation, error) {
	var list []*Reservation
	err := r.db.Preload("Zone").Where("deleted_at IS NULL").Find(&list).Error
	return list, err
}
