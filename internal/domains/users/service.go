package users

import (
	"Spot-Sync/internal/auth"
	"Spot-Sync/internal/domains/users/dto"
	"fmt"
)

type service struct {
	repo       Repository
	jwtService auth.JWTService
}

func NewService(repo Repository, jwtService auth.JWTService) *service {
	return &service{
		repo:       repo,
		jwtService: jwtService,
	}
}

func (s *service) CreateUser(req dto.CreateRequest) (*dto.Response, error) {
	user := User{
		Name:  req.Name,
		Email: req.Email,
		Role:  req.Role,
	}

	if err := user.hashPassword(req.Password); err != nil {
		return nil, err
	}

	err := s.repo.CreateUser(&user)

	if err != nil {
		return nil, err
	}

	response := dto.Response{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
	return &response, nil

}

func (s *service) LoginUser(req dto.LoginRequest) (*dto.LoginResponse, error) {

	user, err := s.repo.GetUserByEmail(req.Email)

	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, ErrorInvalidCredentials
	}

	if err := user.checkPassword(req.Password); err != nil {
		return nil, ErrorInvalidCredentials
	}

	token, err := s.jwtService.GenerateToken(user.ID, req.Email, user.Name, user.Role)
	if err != nil {
		return nil, fmt.Errorf("fail to generate token: %w", err)
	}

	response := dto.LoginResponse{
		Token: token,
		User: dto.LoginUserResponse{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
			Role:  user.Role,
		},
	}

	return &response, nil
}
