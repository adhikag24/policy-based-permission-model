package policies

import (
	"context"
	"log/slog"
	"strings"
)

type Service interface {
	CreatePolicy(ctx context.Context, policy *Policy) (*Policy, error)
	DeletePolicy(ctx context.Context, policyID string) error
	CheckPermission(ctx context.Context, request *CheckPermissionRequest) bool
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) CreatePolicy(ctx context.Context, policy *Policy) (*Policy, error) {
	// If user already has broader policy, no need to add this one.
	// E.g., if user has blogs/*, no need to add blogs/123/*
	if s.isUserHasBroaderPolicy(policy) {
		return nil, ErrUserAlreadyHasBroaderPolicy
	}

	// Delete existing policies that match the resource prefix to avoid duplicates.
	// E.g., if adding blogs/*, remove blogs/123/* first.
	if err := s.repo.DeleteByPrefix(ctx, &DeleteByPrefixRequest{
		AccountID:      policy.AccountID,
		TeamMemberID:   policy.TeamMemberID,
		ResourcePrefix: s.getPrefixByResource(policy.Resource),
		Action:         policy.Action,
	}); err != nil {
		return nil, err
	}

	return s.repo.Create(ctx, policy)
}

func (s *service) isUserHasBroaderPolicy(policy *Policy) bool {
	currentPolicies, err := s.repo.Get(context.Background(), &GetPolicyRequest{
		AccountID:    policy.AccountID,
		TeamMemberID: policy.TeamMemberID,
		Action:       policy.Action,
	})

	if err != nil {
		return false
	}

	return s.hasBroaderPolicy(policy.Resource, currentPolicies)
}

func (s *service) hasBroaderPolicy(resource string, policies []Policy) bool {
	for _, policy := range policies {
		if s.isRootPolicies(policy.Resource) {
			return true
		}
		if strings.HasSuffix(policy.Resource, "/*") {
			prefix := strings.TrimSuffix(policy.Resource, "*")
			if strings.HasPrefix(resource, prefix) {
				return true
			}
		}
	}
	return false
}

func (s *service) isRootPolicies(resource string) bool {
	return resource == "*"
}

func (s *service) getPrefixByResource(resource string) string {
	if resource == "*" {
		return "" // Root access has no prefix.
	}

	if strings.HasSuffix(resource, "/*") {
		return strings.TrimSuffix(resource, "*") // E.g., blogs/* -> blogs/
	}

	return resource + "/"
}

func (s *service) DeletePolicy(ctx context.Context, policyID string) error {
	return s.repo.Delete(ctx, policyID)
}

func (s *service) CheckPermission(ctx context.Context, request *CheckPermissionRequest) bool {
	policies, err := s.repo.Get(ctx, &GetPolicyRequest{
		AccountID:    request.AccountID,
		TeamMemberID: request.TeamMemberID,
		Action:       request.Action,
	})
	if err != nil {
		slog.ErrorContext(ctx, "failed to get policies", "error", err)
		return false
	}

	for _, policy := range policies {
		if s.checkResourceAccess(policy.Resource, request.Resource, request.Action) {
			return true
		}
	}

	return false
}

func (s *service) checkResourceAccess(policyResource, requestResource string, action Action) bool {
	if policyResource == "*" {
		return true // Root access grants all permissions.
	}

	// Exact match for accessing specific resource.
	// E.g., blogs/123/*, blogs
	if policyResource == requestResource {
		return true
	}

	// If user has access to all sub-resources under a resource. E.g., blogs/*
	if strings.HasSuffix(policyResource, "/*") {
		prefix := strings.TrimSuffix(policyResource, "*") // E.g., blogs/* -> blogs/
		return strings.HasPrefix(requestResource, prefix) // E.g., blogs/123 has prefix blogs/
	}

	// If user has access to specific resource and its sub-resources. E.g., blogs/11/pages/12
	// Then they can read blogs/11, but not write.
	if strings.HasPrefix(policyResource, requestResource) && action == ActionRead {
		return true
	}

	return false
}
