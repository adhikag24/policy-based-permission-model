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
	api.GET("/v1/funnels/:id", h.Funnels.GetFunnel)

	api.POST("/v1/blogs/pages", h.Blogs.WriteBlogPage)
	api.GET("/v1/blogs/pages", h.Blogs.ReadBlogPage)
	api.POST("/v1/blogs/settings", h.Blogs.WriteBlogSettings)
	api.GET("/v1/blogs/settings/:id", h.Blogs.ReadBlogSettings)

}
