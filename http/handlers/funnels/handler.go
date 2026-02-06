package handlerfunnels

import (
	"github.com/adhikag24/policy-based-permission-model/domain/funnels"
	"github.com/adhikag24/policy-based-permission-model/http/handlers/shared"
	"github.com/labstack/echo/v5"
)

type Handler struct {
	service funnels.Service
}

func NewHandler(service funnels.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) CreateFunnel(c *echo.Context) error {
	var request CommonRequest[CreateFunnelRequest]
	if err := c.Bind(&request); err != nil {
		return c.JSON(400, shared.Response[any]{
			Code:    400,
			Message: "Invalid request payload",
		})
	}

	err := h.service.CreateFunnel(c.Request().Context(), &funnels.CreateFunnelRequest{
		AccountID:    request.Data.AccountID,
		TeamMemberID: request.Data.TeamMemberID,
		Name:         request.Data.Name,
	})
	if err != nil {
		return c.JSON(500, Response[any]{
			Code: 500,
			Errors: []shared.Errors{
				{
					Code:    "ErrFailedToCreateFunnel",
					Message: "Failed to create funnel",
				},
			},
		})
	}

	return c.JSON(201, Response[any]{
		Code:    201,
		Message: "Successfully created funnel",
	})
}

func (h *Handler) GetFunnel(c *echo.Context) error {
	var request CommonRequest[GetFunnelRequest]
	if err := c.Bind(&request); err != nil {
		return c.JSON(400, shared.Response[any]{
			Code:    400,
			Message: "Invalid request payload",
		})
	}

	funnel, err := h.service.GetFunnel(c.Request().Context(), &funnels.GetFunnelRequest{
		AccountID:    request.Data.AccountID,
		TeamMemberID: request.Data.TeamMemberID,
		FunnelID:     request.Data.FunnelID,
	})
	if err != nil {
		return c.JSON(500, shared.Response[any]{
			Code:    500,
			Message: "Failed to get funnel",
		})
	}

	return c.JSON(200, shared.Response[any]{
		Code:    200,
		Message: "Successfully retrieved funnel",
		Data:    funnel,
	})
}
