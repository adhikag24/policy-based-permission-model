package handlerfunnels

import "github.com/adhikag24/policy-based-permission-model/http/handlers/shared"

type (
	CommonRequest[T any] shared.CommonRequest[T]
	Response[T any]      shared.Response[T]
)

type CreateFunnelRequest struct {
	AccountID    int64  `json:"account_id"`
	TeamMemberID int64  `json:"team_member_id"`
	Name         string `json:"name"`
}

type GetFunnelRequest struct {
	// To simplify put account id and team member id in the request
	AccountID    int64  `json:"account_id"`
	TeamMemberID int64  `json:"team_member_id"`
	FunnelID     string `json:"funnel_id"`
}
