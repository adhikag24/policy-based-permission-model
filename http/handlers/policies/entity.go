package handlerspolicies

import "github.com/adhikag24/policy-based-permission-model/http/handlers/shared"

type CheckPermissionRequest struct {
	AccountID    int64 `json:"account_id"`
	TeamMemberID int64 `json:"team_member_id"`
	// Resource to be accessed.
	Resource string `json:"resource"`
	// Action to be performed.
	Action string `json:"action"`
}

type Policy struct {
	ID           int64  `json:"id"`
	AccountID    int64  `json:"account_id"`
	TeamMemberID int64  `json:"team_member_id"`
	Resource     string `json:"resource"`
	Action       string `json:"action"`
}

type (
	CommonRequest[T any] shared.CommonRequest[T]
	Response[T any]      shared.Response[T]
	Errors               shared.Errors
)
