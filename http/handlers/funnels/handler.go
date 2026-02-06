package handlerfunnels

import (
	"errors"

	"github.com/adhikag24/policy-based-permission-model/domain/blogs"
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
		return h.handleErrorResponse(c, handleErrorResponseSpec{
			err:                     err,
			permissionDeniedCode:    "ErrPermissionDenied",
			permissionDeniedMessage: "Permission denied to create funnel",
			genericErrorCode:        "ErrFailedToCreateFunnel",
			genericErrorMessage:     "Failed to create funnel",
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
		return h.handleErrorResponse(c, handleErrorResponseSpec{
			err:                     err,
			permissionDeniedCode:    "ErrPermissionDenied",
			permissionDeniedMessage: "Permission denied to retrieve funnel",
			genericErrorCode:        "ErrFailedToRetrieveFunnel",
			genericErrorMessage:     "Failed to retrieve funnel",
		})
	}

	return c.JSON(200, shared.Response[any]{
		Code:    200,
		Message: "Successfully retrieved funnel",
		Data:    funnel,
	})
}

type handleErrorResponseSpec struct {
	err                     error
	permissionDeniedCode    string
	permissionDeniedMessage string
	genericErrorCode        string
	genericErrorMessage     string
}

func (h *Handler) handleErrorResponse(c *echo.Context, spec handleErrorResponseSpec) error {
	var (
		err                     = spec.err
		permissionDeniedCode    = spec.permissionDeniedCode
		permissionDeniedMessage = spec.permissionDeniedMessage
		genericErrorCode        = spec.genericErrorCode
		genericErrorMessage     = spec.genericErrorMessage
	)

	if errors.Is(err, blogs.ErrPermissionDenied) {
		return c.JSON(403, Response[any]{
			Code: 403,
			Errors: []shared.Errors{
				{
					Code:    permissionDeniedCode,
					Message: permissionDeniedMessage,
				},
			},
		})
	}
	return c.JSON(500, Response[any]{
		Code: 500,
		Errors: []shared.Errors{
			{
				Code:    genericErrorCode,
				Message: genericErrorMessage,
			},
		},
	})
}
