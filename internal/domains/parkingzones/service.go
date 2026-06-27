package parkings

import (
	"Spot-Sync/internal/domains/parkingzones/dto"
	"time"
)

type service struct {
	repo Repository
}

func NewService(repo Repository) *service {
	return &service{repo: repo}
}

// 3. Create Parking Zone (POST /api/v1/zones)
func (s *service) CreateZone(req dto.CreateZoneRequest) (*dto.ZoneResponse, error) {
	zone := Zone{
		Name:          req.Name,
		Type:          ZoneType(req.Type),
		TotalCapacity: req.TotalCapacity,
		PricePerHour:  req.PricePerHour,
	}

	err := s.repo.Create(&zone)
	if err != nil {
		return nil, err
	}

	return &dto.ZoneResponse{
		ID:             zone.ID,
		Name:           zone.Name,
		Type:           string(zone.Type),
		TotalCapacity:  zone.TotalCapacity,
		PricePerHour:   zone.PricePerHour,
		CreatedAt:      zone.CreatedAt.Format(time.RFC3339),
		UpdatedAt:      zone.UpdatedAt.Format(time.RFC3339),
		AvailableSpots: nil, // Automatically omitted from JSON output
	}, nil
}

// 4. Get All Parking Zones (GET /api/v1/zones)
func (s *service) GetZones() ([]*dto.ZoneResponse, error) {
	zones, err := s.repo.GetAll()
	if err != nil {
		return nil, err
	}

	var responses []*dto.ZoneResponse
	for _, z := range zones {
		spotsCopy := z.AvailableSpots // Addressable local copy
		responses = append(responses, &dto.ZoneResponse{
			ID:             z.ID,
			Name:           z.Name,
			Type:           string(z.Type),
			TotalCapacity:  z.TotalCapacity,
			AvailableSpots: &spotsCopy, // Populates available_spots field
			PricePerHour:   z.PricePerHour,
			CreatedAt:      z.CreatedAt.Format(time.RFC3339),
			UpdatedAt:      "", // Automatically omitted from JSON output
		})
	}
	return responses, nil
}

// 5. Get Single Parking Zone (GET /api/v1/zones/:id)
func (s *service) GetZoneByID(id uint) (*dto.ZoneResponse, error) {
	z, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	spotsCopy := z.AvailableSpots
	return &dto.ZoneResponse{
		ID:             z.ID,
		Name:           z.Name,
		Type:           string(z.Type),
		TotalCapacity:  z.TotalCapacity,
		AvailableSpots: &spotsCopy, // Populates available_spots field
		PricePerHour:   z.PricePerHour,
		CreatedAt:      z.CreatedAt.Format(time.RFC3339),
		UpdatedAt:      "", // Automatically omitted from JSON output
	}, nil
}
