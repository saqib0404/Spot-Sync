package users_test

import (
	"Spot-Sync/internal/domains/users/dto"
	"testing"

	"github.com/go-playground/validator/v10"
)

func TestCreateRequestValidation(t *testing.T) {
	validate := validator.New()

	tests := []struct {
		name    string
		req     dto.CreateRequest
		wantErr bool
	}{
		{
			name: "valid admin role",
			req: dto.CreateRequest{
				Name:     "Admin User",
				Email:    "admin@example.com",
				Password: "password123",
				Role:     "admin",
			},
			wantErr: false,
		},
		{
			name: "valid driver role",
			req: dto.CreateRequest{
				Name:     "Driver User",
				Email:    "driver@example.com",
				Password: "password123",
				Role:     "driver",
			},
			wantErr: false,
		},
		{
			name: "invalid role",
			req: dto.CreateRequest{
				Name:     "Invalid User",
				Email:    "invalid@example.com",
				Password: "password123",
				Role:     "passenger",
			},
			wantErr: true,
		},
		{
			name: "empty role",
			req: dto.CreateRequest{
				Name:     "Invalid User",
				Email:    "invalid@example.com",
				Password: "password123",
				Role:     "",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validate.Struct(&tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("validate.Struct() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
