package handlerfunnels

import (
	"errors"
	"strconv"

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

	// Mandatory headers.
	accountID, teamMemberID, err := h.getMandatoryHeaders(c)
	if err != nil {
		return c.JSON(400, shared.Response[any]{
			Code:    400,
			Message: "Missing mandatory headers",
		})
	}

	err = h.service.CreateFunnel(c.Request().Context(), &funnels.CreateFunnelRequest{
		AccountID:    accountID,
		TeamMemberID: teamMemberID,
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
	funnelID := c.Param("id")

	// Mandatory headers.
	accountID, teamMemberID, err := h.getMandatoryHeaders(c)
	if err != nil {
		return c.JSON(400, shared.Response[any]{
			Code:    400,
			Message: "Missing mandatory headers",
		})
	}

	funnel, err := h.service.GetFunnel(c.Request().Context(), &funnels.GetFunnelRequest{
		AccountID:    accountID,
		TeamMemberID: teamMemberID,
		FunnelID:     funnelID,
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

	if errors.Is(err, funnels.ErrPermissionDenied) {
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

func (h *Handler) getMandatoryHeaders(c *echo.Context) (accountID int64, teamMemberID int64, err error) {
	accountIDStr := c.Request().Header.Get("X-Account-ID")
	teamMemberIDStr := c.Request().Header.Get("X-Team-Member-ID")

	if accountIDStr == "" || teamMemberIDStr == "" {
		return 0, 0, errors.New("missing mandatory headers")
	}

	accountIDInt, err := strconv.Atoi(accountIDStr)
	if err != nil {
		return 0, 0, errors.New("invalid X-Account-ID header")
	}

	teamMemberIDInt, err := strconv.Atoi(teamMemberIDStr)
	if err != nil {
		return 0, 0, errors.New("invalid X-Team-Member-ID header")
	}

	return int64(accountIDInt), int64(teamMemberIDInt), nil
}
