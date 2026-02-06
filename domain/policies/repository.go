package policies

import "context"

type Repository interface {
	Create(ctx context.Context, policy *Policy) (*Policy, error)
	Delete(ctx context.Context, policyID string) error
	Get(ctx context.Context, request *GetPolicyRequest) ([]Policy, error)
	DeleteByPrefix(ctx context.Context, request *DeleteByPrefixRequest) error
}

// Retreive policy based on AccountID, TeamMemberID, and Action.
type GetPolicyRequest struct {
	AccountID    int64
	TeamMemberID int64
	Action       Action
}

type DeleteByPrefixRequest struct {
	AccountID      int64
	TeamMemberID   int64
	ResourcePrefix string
	Action         Action
}
