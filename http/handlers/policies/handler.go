package handlerspolicies

import (
	"errors"

	"github.com/adhikag24/policy-based-permission-model/domain/policies"
	"github.com/adhikag24/policy-based-permission-model/http/handlers/shared"
	"github.com/labstack/echo/v5"
)

type Handler struct {
	service policies.Service
}

func NewHandler(service policies.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) CreatePolicy(c *echo.Context) error {
	var request CommonRequest[Policy]
	if err := c.Bind(&request); err != nil {
		return c.JSON(400, shared.Response[any]{
			Code:    400,
			Message: "Invalid request payload",
		})
	}

	requestContext := c.Request().Context()
	policyDomainRequest := policies.Policy{
		AccountID:    request.Data.AccountID,
		TeamMemberID: request.Data.TeamMemberID,
		Resource:     request.Data.Resource,
		Action:       policies.Action(request.Data.Action),
	}
	policy, err := h.service.CreatePolicy(requestContext, &policyDomainRequest)
	if err != nil {
		if errors.Is(err, policies.ErrUserAlreadyHasBroaderPolicy) {
			// Treat as success response.
			return c.JSON(201, Response[any]{
				Code:    201,
				Message: "Successfully created policy",
			})
		}
		// Generic error response.
		return c.JSON(500, Response[any]{
			Code: 500,
			Errors: []shared.Errors{
				{
					Code:    "ErrFailedToCreatePolicy",
					Message: "Failed to create policy",
				},
			},
		})
	}

	responsePolicy := &Policy{
		ID:           policy.ID,
		AccountID:    policy.AccountID,
		TeamMemberID: policy.TeamMemberID,
		Resource:     policy.Resource,
		Action:       string(policy.Action),
	}

	return c.JSON(201, Response[*Policy]{
		Code:    201,
		Message: "Successfully created policy",
		Data:    responsePolicy,
	})
}

func (h *Handler) DeletePolicy(c *echo.Context) error {
	policyID := c.Param("id")
	if policyID == "" {
		return c.JSON(400, Response[any]{
			Code: 400,
			Errors: []shared.Errors{
				{
					Code:    "ErrPolicyIDRequired",
					Message: "Policy ID is required",
				},
			},
		})
	}

	requestContext := c.Request().Context()
	err := h.service.DeletePolicy(requestContext, policyID)
	if err != nil {
		return c.JSON(500, Response[any]{
			Code: 500,
			Errors: []shared.Errors{
				{
					Code:    "ErrFailedToDeletePolicy",
					Message: "Failed to delete policy",
				},
			},
		})
	}

	return c.JSON(200, Response[any]{
		Code:    200,
		Message: "Successfully deleted policy",
	})
}

func (h *Handler) CheckPermission(c *echo.Context) error {
	var request CommonRequest[CheckPermissionRequest]
	if err := c.Bind(&request); err != nil {
		return c.JSON(400, Response[any]{
			Code:    400,
			Message: "Invalid request payload",
		})
	}

	requestContext := c.Request().Context()
	isPermitted := h.service.CheckPermission(requestContext, &policies.CheckPermissionRequest{
		AccountID:    request.Data.AccountID,
		TeamMemberID: request.Data.TeamMemberID,
		Resource:     request.Data.Resource,
		Action:       policies.Action(request.Data.Action),
	})
	if !isPermitted {
		return c.JSON(403, Response[any]{
			Code: 403,
			Errors: []shared.Errors{
				{
					Code:    "ErrPermissionDenied",
					Message: "Permission denied",
				},
			},
		})
	}

	return c.JSON(200, Response[any]{
		Code:    200,
		Message: "Permission is valid",
	})
}
