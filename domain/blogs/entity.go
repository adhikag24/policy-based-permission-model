package blogs

type WriteBlogPageRequest struct {
	AccountID    int64
	TeamMemberID int64
	Title        string
	Content      string
	PageID       int64
}

type WriteBlogSettingsRequest struct {
	AccountID    int64
	TeamMemberID int64
	BlogID       int64
	Title        string
	Content      string
}

type ReadBlogPageRequest struct {
	AccountID    int64
	TeamMemberID int64
	BlogID       int64
	PageID       int64
}

type ReadBlogSettingsRequest struct {
	AccountID    int64
	TeamMemberID int64
	BlogID       int64
}
