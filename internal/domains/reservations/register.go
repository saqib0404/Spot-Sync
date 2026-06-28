package reservations

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
	service := NewService(repo)
	jwtService := auth.NewJWTService(cfg.JWTSecret)
	handler := NewHandler(service)

	authMid := middleware.AuthMiddleware(jwtService)

	adminOnly := func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			role, ok := c.Get("user_role").(string)
			if !ok || role != "admin" {
				return c.JSON(http.StatusForbidden, map[string]string{"error": "Admin access required"})
			}
			return next(c)
		}
	}

	// Router Group
	api := e.Group("/api/v1/reservations", authMid)

	api.POST("", handler.Create)
	api.GET("/my-reservations", handler.GetMy)
	api.DELETE("/:id", handler.Cancel)
	api.GET("", handler.GetAllAdmin, adminOnly)
}
