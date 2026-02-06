package blogs

import "github.com/adhikag24/policy-based-permission-model/http/handlers/shared"

type (
	CommonRequest[T any] shared.CommonRequest[T]
	Response[T any]      shared.Response[T]
)

type WriteBlogPageRequest struct {
	AccountID    int64  `json:"account_id"`
	TeamMemberID int64  `json:"team_member_id"`
	PageID       int64  `json:"page_id"`
	Title        string `json:"title"`
	Content      string `json:"content"`
}

type ReadBlogPageRequest struct {
	AccountID    int64 `json:"account_id"`
	TeamMemberID int64 `json:"team_member_id"`
	PageID       int64 `json:"page_id"`
}

type ReadBlogSettingsRequest struct {
	AccountID    int64 `json:"account_id"`
	TeamMemberID int64 `json:"team_member_id"`
	BlogID       int64 `json:"blog_id"`
}

type WriteBlogSettingsRequest struct {
	AccountID    int64  `json:"account_id"`
	TeamMemberID int64  `json:"team_member_id"`
	BlogID       int64  `json:"blog_id"`
	Title        string `json:"title"`
	Content      string `json:"content"`
}
