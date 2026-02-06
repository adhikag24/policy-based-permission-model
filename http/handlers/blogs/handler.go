package blogs

import (
	"errors"

	"github.com/adhikag24/policy-based-permission-model/domain/blogs"
	"github.com/adhikag24/policy-based-permission-model/http/handlers/shared"
	"github.com/labstack/echo/v5"
)

type Handler struct {
	blogsService blogs.Service
}

func NewHandler(blogsService blogs.Service) *Handler {
	return &Handler{blogsService: blogsService}
}

func (h *Handler) WriteBlogPage(c *echo.Context) error {
	var request CommonRequest[WriteBlogPageRequest]
	if err := c.Bind(&request); err != nil {
		return err
	}

	err := h.blogsService.WriteBlogPage(c.Request().Context(), &blogs.WriteBlogPageRequest{
		AccountID:    request.Data.AccountID,
		TeamMemberID: request.Data.TeamMemberID,
		PageID:       request.Data.PageID,
		Content:      request.Data.Content,
	})
	if err != nil {
		return h.handleErrorResponse(c, handleErrorResponseSpec{
			err:                     err,
			permissionDeniedCode:    "ErrPermissionDenied",
			permissionDeniedMessage: "Permission denied to write blog page",
			genericErrorCode:        "ErrFailedToWriteBlogPage",
			genericErrorMessage:     "Failed to write blog page",
		})
	}

	return c.JSON(201, Response[any]{
		Code:    201,
		Message: "Successfully wrote blog page",
	})
}

func (h *Handler) ReadBlogPage(c *echo.Context) error {
	var request CommonRequest[ReadBlogPageRequest]
	if err := c.Bind(&request); err != nil {
		return err
	}

	err := h.blogsService.ReadBlogPage(c.Request().Context(), &blogs.ReadBlogPageRequest{
		AccountID:    request.Data.AccountID,
		TeamMemberID: request.Data.TeamMemberID,
		PageID:       request.Data.PageID,
	})
	if err != nil {
		return h.handleErrorResponse(c, handleErrorResponseSpec{
			err:                     err,
			permissionDeniedCode:    "ErrPermissionDenied",
			permissionDeniedMessage: "Permission denied to read blog page",
			genericErrorCode:        "ErrFailedToReadBlogPage",
			genericErrorMessage:     "Failed to read blog page",
		})
	}

	return c.JSON(200, Response[any]{
		Code:    200,
		Message: "Successfully read blog page",
	})
}

func (h *Handler) ReadBlogSettings(c *echo.Context) error {
	var request CommonRequest[ReadBlogSettingsRequest]
	if err := c.Bind(&request); err != nil {
		return err
	}

	err := h.blogsService.ReadBlogSettings(c.Request().Context(), &blogs.ReadBlogSettingsRequest{
		AccountID:    request.Data.AccountID,
		TeamMemberID: request.Data.TeamMemberID,
	})
	if err != nil {
		return h.handleErrorResponse(c, handleErrorResponseSpec{
			err:                     err,
			permissionDeniedCode:    "ErrPermissionDenied",
			permissionDeniedMessage: "Permission denied to read blog settings",
			genericErrorCode:        "ErrFailedToReadBlogSettings",
			genericErrorMessage:     "Failed to read blog settings",
		})
	}

	return c.JSON(200, Response[any]{
		Code:    200,
		Message: "Successfully read blog settings",
	})
}

func (h *Handler) WriteBlogSettings(c *echo.Context) error {
	var request CommonRequest[WriteBlogSettingsRequest]
	if err := c.Bind(&request); err != nil {
		return err
	}

	err := h.blogsService.WriteBlogSettings(c.Request().Context(), &blogs.WriteBlogSettingsRequest{
		AccountID:    request.Data.AccountID,
		TeamMemberID: request.Data.TeamMemberID,
		BlogID:       request.Data.BlogID,
		Title:        request.Data.Title,
		Content:      request.Data.Content,
	})
	if err != nil {
		return c.JSON(500, Response[any]{
			Code: 500,
			Errors: []shared.Errors{
				{
					Code:    "ErrFailedToWriteBlogSettings",
					Message: "Failed to write blog settings",
				},
			},
		})
	}

	return c.JSON(201, Response[any]{
		Code:    201,
		Message: "Successfully wrote blog settings",
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
