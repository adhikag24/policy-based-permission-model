package policies

type Action string

const (
	ActionRead  Action = "read"
	ActionWrite Action = "write"
)

type Policy struct {
	ID           int64
	AccountID    int64
	TeamMemberID int64
	Resource     string
	Action       Action
}

type CheckPermissionRequest struct {
	AccountID    int64
	TeamMemberID int64
	Resource     string // E.g., blogs
	Action       Action
}
