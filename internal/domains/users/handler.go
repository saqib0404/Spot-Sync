package users

import (
	"Spot-Sync/internal/domains/users/dto"
	"Spot-Sync/internal/httpresponse"
	"errors"
	"net/http"

	"github.com/labstack/echo/v5"
)

type handler struct {
	service *service
}

func NewHandler(service *service) *handler {
	return &handler{
		service: service,
	}
}

func (h *handler) CreateUser(c *echo.Context) error {
	var req dto.CreateRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid request payload",
			Details: err.Error(),
		})
	}

	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Validation failed",
			Details: err.Error(),
		})
	}

	response, err := h.service.CreateUser(req)
	if err != nil {

		if errors.Is(err, ErrorUserAlreadyExists) {
			return c.JSON(http.StatusConflict, httpresponse.ErrorResponse{
				Code:    http.StatusConflict,
				Message: "User with this email already exists",
				Details: err.Error(),
			})
		}

		return c.JSON(http.StatusInternalServerError, httpresponse.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to create user",
			Details: err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, map[string]any{
		"success": true,
		"message": "User registered successfully",
		"data":    response,
	})
}

// func (h *handler) LoginUser(c *echo.Context) error {
// 	var req dto.LoginRequest

// 	if err := c.Bind(&req); err != nil {
// 		return c.JSON(http.StatusBadRequest, httpresponse.ErrorResponse{
// 			Code:    http.StatusBadRequest,
// 			Message: "Invalid request payload",
// 			Details: err.Error(),
// 		})
// 	}

// 	if err := c.Validate(&req); err != nil {
// 		return c.JSON(http.StatusBadRequest, httpresponse.ErrorResponse{
// 			Code:    http.StatusBadRequest,
// 			Message: "Validation failed",
// 			Details: err.Error(),
// 		})
// 	}

// 	response, err := h.service.LoginUser(req)
// 	if err != nil {

// 		if errors.Is(err, ErrorInvalidCredentials) {
// 			return c.JSON(http.StatusUnauthorized, httpresponse.ErrorResponse{
// 				Code:    http.StatusUnauthorized,
// 				Message: "Invalid email or password",
// 				Details: err.Error(),
// 			})
// 		}

// 		return c.JSON(http.StatusInternalServerError, httpresponse.ErrorResponse{
// 			Code:    http.StatusInternalServerError,
// 			Message: "Failed to login user",
// 			Details: err.Error(),
// 		})
// 	}

// 	return c.JSON(http.StatusOK, response)
// }

// func (h *handler) GetMe(c *echo.Context) error {
// 	userID, ok := c.Get("user_id").(uint)
// 	if !ok {
// 		return c.JSON(http.StatusUnauthorized, httpresponse.ErrorResponse{
// 			Code:    http.StatusUnauthorized,
// 			Message: "Cannot get user info",
// 			Details: "Missing user id in context",
// 		})
// 	}

// 	email, _ := c.Get("user_email").(string)
// 	name, _ := c.Get("user_name").(string)

// 	return c.JSON(http.StatusOK, dto.Response{
// 		ID:    userID,
// 		Name:  name,
// 		Email: email,
// 	})
// }
