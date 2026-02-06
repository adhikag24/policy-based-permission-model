package mysqlpolicies

import (
	"time"

	"github.com/adhikag24/policy-based-permission-model/domain/policies"
)

type PolicyModel struct {
	ID           int64 `gorm:"primaryKey"`
	AccountID    int64
	TeamMemberID int64
	Resource     string
	Action       string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func (PolicyModel) TableName() string {
	return "policies"
}

func ToDomain(m PolicyModel) policies.Policy {
	return policies.Policy{
		ID:           m.ID,
		AccountID:    m.AccountID,
		TeamMemberID: m.TeamMemberID,
		Resource:     m.Resource,
		Action:       policies.Action(m.Action),
	}
}

func FromDomain(p policies.Policy) PolicyModel {
	return PolicyModel{
		ID:           p.ID,
		AccountID:    p.AccountID,
		TeamMemberID: p.TeamMemberID,
		Resource:     p.Resource,
		Action:       string(p.Action),
	}
}
