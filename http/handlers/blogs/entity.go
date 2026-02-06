package blogs

import "github.com/adhikag24/policy-based-permission-model/http/handlers/shared"

type (
	CommonRequest[T any] shared.CommonRequest[T]
	Response[T any]      shared.Response[T]
)

type WriteBlogPageRequest struct {
	PageID  string `json:"page_id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

type WriteBlogSettingsRequest struct {
	BlogID  string `json:"blog_id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}
