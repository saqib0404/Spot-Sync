package server

import (
	"Spot-Sync/internal/config"
	parkings "Spot-Sync/internal/domains/parkingZones"
	"Spot-Sync/internal/domains/users"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
)

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i any) error {
	if err := cv.validator.Struct(i); err != nil {
		// Optionally, you could return the error to give each route more control over the status code
		return echo.ErrBadRequest.Wrap(err)
	}
	return nil
}

func Start(cfg *config.Config, db *gorm.DB) {
	// db.AutoMigrate(&user.User{}, &event.Event{}, &booking.Booking{})
	db.AutoMigrate(&users.User{}, &parkings.Zone{})

	e := echo.New()
	// e.Use(middleware.RequestLogger())
	e.Validator = &CustomValidator{validator: validator.New()}

	e.GET("/", func(c *echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	users.RegisterRoutes(e, db, cfg)
	parkings.RegisterRoutes(e, db, cfg)
	// booking.RegisterRoutes(e, db, cfg)

	port := fmt.Sprintf(":%s", cfg.Port)
	if err := e.Start(port); err != nil {
		e.Logger.Error("failed to start server", "error", err)
	}
}
