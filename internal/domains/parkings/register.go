package parkings

import (
	"Spot-Sync/internal/auth"
	"Spot-Sync/internal/config"
	middleware "Spot-Sync/internal/middlewares"
	"net/http"

	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
)

func RegisterRoutes(e *echo.Echo, db *gorm.DB, cfg *config.Config) {
	repo := NewRepository(db)
	jwtService := auth.NewJWTService(cfg.JWTSecret)
	service := NewService(repo)
	handler := NewHandler(service)

	// Auth Middleware instance
	authMid := middleware.AuthMiddleware(jwtService)

	// Admin authorization checkpoint
	adminOnly := func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			role, ok := c.Get("user_role").(string)
			if !ok || role != "admin" {
				return c.JSON(http.StatusForbidden, map[string]string{
					"error": "Access denied: Admins only",
				})
			}
			return next(c)
		}
	}

	// Base API group paths
	api := e.Group("/api/v1/zones")

	// Routes
	api.POST("", handler.CreateZone, authMid, adminOnly) // Admin Only
	api.GET("", handler.GetZones)                        // Public
	api.GET("/:id", handler.GetZoneByID)                 // Public
}
