package reservations

import (
	"Spot-Sync/internal/domains/reservations/dto"
	"errors"
	"time"
)

type service struct {
	repo Repository
}

func NewService(repo Repository) *service {
	return &service{repo: repo}
}

func (s *service) MakeReservation(userID uint, req dto.CreateReservationRequest) (*dto.ReservationResponse, error) {
	res, err := s.repo.CreateAtomic(userID, req.ZoneID, req.LicensePlate)
	if err != nil {
		return nil, err
	}

	return &dto.ReservationResponse{
		ID:           res.ID,
		UserID:       res.UserID,
		ZoneID:       res.ZoneID,
		LicensePlate: res.LicensePlate,
		Status:       string(res.Status),
		CreatedAt:    res.CreatedAt.Format(time.RFC3339),
		UpdatedAt:    res.UpdatedAt.Format(time.RFC3339),
	}, nil
}

func (s *service) GetUserReservations(userID uint) ([]*dto.ReservationResponse, error) {
	list, err := s.repo.GetByUserID(userID)
	if err != nil {
		return nil, err
	}

	var mapped []*dto.ReservationResponse
	for _, item := range list {
		mapped = append(mapped, &dto.ReservationResponse{
			ID:           item.ID,
			LicensePlate: item.LicensePlate,
			Status:       string(item.Status),
			CreatedAt:    item.CreatedAt.Format(time.RFC3339),
			Zone: &dto.NestedZoneInfo{
				ID:   item.Zone.ID,
				Name: item.Zone.Name,
				Type: string(item.Zone.Type),
			},
		})
	}
	return mapped, nil
}

func (s *service) CancelReservation(userID uint, role string, resID uint) error {
	res, err := s.repo.GetByID(resID)
	if err != nil {
		return err
	}

	// Route Protection Guardrail: Drivers can only drop their own items
	if role != "admin" && res.UserID != userID {
		return errors.New("forbidden")
	}

	if res.Status == StatusCancelled {
		return ErrReservationAlreadyCancelled
	}

	return s.repo.UpdateStatus(res, StatusCancelled)
}

func (s *service) GetAllReservationsAdmin() ([]*dto.ReservationResponse, error) {
	list, err := s.repo.GetAllAdmin()
	if err != nil {
		return nil, err
	}

	var mapped []*dto.ReservationResponse
	for _, item := range list {
		mapped = append(mapped, &dto.ReservationResponse{
			ID:           item.ID,
			UserID:       item.UserID,
			ZoneID:       item.ZoneID,
			LicensePlate: item.LicensePlate,
			Status:       string(item.Status),
			CreatedAt:    item.CreatedAt.Format(time.RFC3339),
			Zone: &dto.NestedZoneInfo{
				ID:   item.Zone.ID,
				Name: item.Zone.Name,
				Type: string(item.Zone.Type),
			},
		})
	}
	return mapped, nil
}
