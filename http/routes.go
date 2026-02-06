package http

import (
	"github.com/labstack/echo/v5"
)

func RegisterRoutes(e *echo.Echo, h *Handlers) {
	api := e.Group("/api")

	api.POST("/v1/policies", h.Policies.CreatePolicy)
	api.DELETE("/v1/policies/:id", h.Policies.DeletePolicy)
	api.POST("/v1/policies/check-permission", h.Policies.CheckPermission)

	api.POST("/v1/funnels", h.Funnels.CreateFunnel)
	api.POST("/v1/funnels/get", h.Funnels.GetFunnel)
}
