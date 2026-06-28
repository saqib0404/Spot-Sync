package reservations

import (
	"Spot-Sync/internal/domains/reservations/dto"
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

func resErrorResponse(c *echo.Context, err error) error {
	if errors.Is(err, ErrZoneFull) {
		return c.JSON(http.StatusConflict, map[string]string{"error": err.Error()})
	}
	if errors.Is(err, ErrReservationNotFound) {
		return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
	}
	if err.Error() == "forbidden" {
		return c.JSON(http.StatusForbidden, map[string]string{"error": "Access denied: cannot modify this resource"})
	}
	return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
}

func (h *handler) Create(c *echo.Context) error {
	uid := c.Get("user_id").(uint)
	var req dto.CreateReservationRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
	}

	res, err := h.service.MakeReservation(uid, req)
	if err != nil {
		return resErrorResponse(c, err)
	}

	return c.JSON(http.StatusCreated, dto.ReservationAPIEnvelope{
		Success: true,
		Message: "Reservation confirmed successfully",
		Data:    res,
	})
}

func (h *handler) GetMy(c *echo.Context) error {
	uid := c.Get("user_id").(uint)
	data, err := h.service.GetUserReservations(uid)
	if err != nil {
		return resErrorResponse(c, err)
	}

	return c.JSON(http.StatusOK, dto.ReservationAPIEnvelope{
		Success: true,
		Message: "My reservations retrieved successfully",
		Data:    data,
	})
}

func (h *handler) Cancel(c *echo.Context) error {
	uid := c.Get("user_id").(uint)
	role := c.Get("user_role").(string)
	id, _ := strconv.Atoi(c.Param("id"))

	if err := h.service.CancelReservation(uid, role, uint(id)); err != nil {
		return resErrorResponse(c, err)
	}

	return c.JSON(http.StatusOK, map[string]any{
		"success": true,
		"message": "Reservation cancelled successfully",
	})
}

func (h *handler) GetAllAdmin(c *echo.Context) error {
	data, err := h.service.GetAllReservationsAdmin()
	if err != nil {
		return resErrorResponse(c, err)
	}

	return c.JSON(http.StatusOK, dto.ReservationAPIEnvelope{
		Success: true,
		Message: "All reservations retrieved successfully",
		Data:    data,
	})
}
