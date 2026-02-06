package mysqlpolicies

import (
	"context"
	"fmt"

	"github.com/adhikag24/policy-based-permission-model/domain/policies"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, policy *policies.Policy) (*policies.Policy, error) {
	policyModel := FromDomain(*policy)
	if err := r.db.WithContext(ctx).Create(&policyModel).Error; err != nil {
		return nil, err
	}
	response := ToDomain(policyModel)
	return &response, nil
}

func (r *Repository) Delete(ctx context.Context, policyID string) error {
	if err := r.db.WithContext(ctx).Delete(&PolicyModel{}, policyID).Error; err != nil {
		return err
	}
	return nil
}

// Retreives list of policies based on account ID, team member ID, and action.
func (r *Repository) Get(ctx context.Context, request *policies.GetPolicyRequest) ([]policies.Policy, error) {
	var policyModels []PolicyModel
	err := r.db.WithContext(ctx).Where("account_id = ? AND team_member_id = ? AND action = ?",
		request.AccountID, request.TeamMemberID, string(request.Action)).Find(&policyModels).Error
	if err != nil {
		return nil, err
	}
	var policies []policies.Policy
	for _, pm := range policyModels {
		policies = append(policies, ToDomain(pm))
	}
	return policies, nil
}

func (r *Repository) DeleteByPrefix(ctx context.Context, request *policies.DeleteByPrefixRequest) error {
	prefixLike := fmt.Sprintf("%s%%", request.ResourcePrefix)
	err := r.db.WithContext(ctx).Where("account_id = ? AND team_member_id = ? AND resource LIKE ? AND action = ?",
		request.AccountID, request.TeamMemberID, prefixLike, string(request.Action)).Delete(&PolicyModel{}).Error
	if err != nil {
		return err
	}
	return nil
}
