package blogs

import (
	"errors"
	"strconv"

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

	accountID, teamMemberID, err := h.getMandatoryHeaders(c)
	if err != nil {
		return c.JSON(400, shared.Response[any]{
			Code: 400,
			Errors: []shared.Errors{
				{
					Code:    "ErrMissingMandatoryHeaders",
					Message: "Missing mandatory headers",
				},
			},
		})
	}

	err = h.blogsService.WriteBlogPage(c.Request().Context(), &blogs.WriteBlogPageRequest{
		AccountID:    accountID,
		TeamMemberID: teamMemberID,
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
	pageID := c.QueryParam("page_id")
	blogID := c.QueryParam("blog_id")

	accountID, teamMemberID, err := h.getMandatoryHeaders(c)
	if err != nil {
		return c.JSON(400, shared.Response[any]{
			Code: 400,
			Errors: []shared.Errors{
				{
					Code:    "ErrMissingMandatoryHeaders",
					Message: "Missing mandatory headers",
				},
			},
		})
	}

	err = h.blogsService.ReadBlogPage(c.Request().Context(), &blogs.ReadBlogPageRequest{
		AccountID:    accountID,
		TeamMemberID: teamMemberID,
		PageID:       pageID,
		BlogID:       blogID,
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
	blogID := c.Param("id")

	accountID, teamMemberID, err := h.getMandatoryHeaders(c)
	if err != nil {
		return c.JSON(400, shared.Response[any]{
			Code: 400,
			Errors: []shared.Errors{
				{
					Code:    "ErrMissingMandatoryHeaders",
					Message: "Missing mandatory headers",
				},
			},
		})
	}

	err = h.blogsService.ReadBlogSettings(c.Request().Context(), &blogs.ReadBlogSettingsRequest{
		AccountID:    accountID,
		TeamMemberID: teamMemberID,
		BlogID:       blogID,
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

	accountID, teamMemberID, err := h.getMandatoryHeaders(c)
	if err != nil {
		return h.handleErrorResponse(c, handleErrorResponseSpec{
			err:                     err,
			permissionDeniedCode:    "ErrPermissionDenied",
			permissionDeniedMessage: "Permission denied to write blog settings",
			genericErrorCode:        "ErrFailedToWriteBlogSettings",
			genericErrorMessage:     "Failed to write blog settings",
		})
	}

	err = h.blogsService.WriteBlogSettings(c.Request().Context(), &blogs.WriteBlogSettingsRequest{
		AccountID:    accountID,
		TeamMemberID: teamMemberID,
		BlogID:       request.Data.BlogID,
		Title:        request.Data.Title,
		Content:      request.Data.Content,
	})
	if err != nil {
		return h.handleErrorResponse(c, handleErrorResponseSpec{
			err:                     err,
			permissionDeniedCode:    "ErrPermissionDenied",
			permissionDeniedMessage: "Permission denied to write blog settings",
			genericErrorCode:        "ErrFailedToWriteBlogSettings",
			genericErrorMessage:     "Failed to write blog settings",
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
