package users

import (
	"Spot-Sync/internal/auth"
	"Spot-Sync/internal/config"

	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
)

func RegisterRoutes(e *echo.Echo, db *gorm.DB, cfg *config.Config) {
	userRepository := NewRepository(db)
	jwtService := auth.NewJWTService(cfg.JWTSecret)
	userService := NewService(userRepository, jwtService)
	userhandler := NewHandler(userService)

	api := e.Group("/api/v1/auth")

	api.POST("/register", userhandler.CreateUser) //api/v1/auth/register
	api.POST("/login", userhandler.LoginUser)     //api/v1/auth/login
	// api.GET("/get-me", userhandler.GetMe, middleware.AuthMiddleware(jwtService))
}
