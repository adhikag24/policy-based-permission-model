package funnels

type Funnel struct {
	FunnelID string
	Name     string
}

type CreateFunnelRequest struct {
	AccountID    int64
	TeamMemberID int64
	Name         string
}

type GetFunnelRequest struct {
	AccountID    int64
	TeamMemberID int64
	FunnelID     string
}

type EditFunnelRequest struct {
	AccountID    int64
	TeamMemberID int64
	FunnelID     string
}
