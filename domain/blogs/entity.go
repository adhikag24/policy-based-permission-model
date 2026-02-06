package blogs

type WriteBlogPageRequest struct {
	AccountID    int64
	TeamMemberID int64
	Title        string
	Content      string
	PageID       string
}

type WriteBlogSettingsRequest struct {
	AccountID    int64
	TeamMemberID int64
	BlogID       string
	Title        string
	Content      string
}

type ReadBlogPageRequest struct {
	AccountID    int64
	TeamMemberID int64
	BlogID       string
	PageID       string
}

type ReadBlogSettingsRequest struct {
	AccountID    int64
	TeamMemberID int64
	BlogID       string
}
