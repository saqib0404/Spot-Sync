package parkings

import (
	"Spot-Sync/internal/domains/parkings/dto"
	"errors"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v5"
)

type handler struct {
	service *service
}

func NewHandler(service *service) *handler {
	return &handler{service: service}
}

type GlobalSuccessResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

type ErrorPayload struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

func zoneErrorResponse(c *echo.Context, err error) error {
	if errors.Is(err, ErrZoneNotFound) {
		return c.JSON(http.StatusNotFound, ErrorPayload{
			Code:    http.StatusNotFound,
			Message: "Parking zone not found",
		})
	}

	return c.JSON(http.StatusInternalServerError, ErrorPayload{
		Code:    http.StatusInternalServerError,
		Message: "Something went wrong",
		Details: err.Error(),
	})
}

func (h *handler) CreateZone(c *echo.Context) error {
	var req dto.CreateZoneRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorPayload{
			Code:    http.StatusBadRequest,
			Message: "Invalid request payload",
			Details: err.Error(),
		})
	}

	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorPayload{
			Code:    http.StatusBadRequest,
			Message: "Validation failed",
			Details: err.Error(),
		})
	}

	response, err := h.service.CreateZone(req)
	if err != nil {
		return zoneErrorResponse(c, err)
	}

	return c.JSON(http.StatusCreated, GlobalSuccessResponse{
		Success: true,
		Message: "Parking zone created successfully",
		Data:    response,
	})
}

func (h *handler) GetZones(c *echo.Context) error {
	zones, err := h.service.GetZones()
	if err != nil {
		return zoneErrorResponse(c, err)
	}

	return c.JSON(http.StatusOK, GlobalSuccessResponse{
		Success: true,
		Message: "Parking zones retrieved successfully",
		Data:    zones,
	})
}

func (h *handler) GetZoneByID(c *echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorPayload{
			Code:    http.StatusBadRequest,
			Message: "Invalid zone ID",
			Details: err.Error(),
		})
	}

	zone, err := h.service.GetZoneByID(uint(id))
	if err != nil {
		return zoneErrorResponse(c, err)
	}

	return c.JSON(http.StatusOK, GlobalSuccessResponse{
		Success: true,
		Message: "Parking zone retrieved successfully",
		Data:    zone,
	})
}
