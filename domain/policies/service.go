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
	// If user already has broader policy, reject lower level policy.
	// E.g., if user has blogs/* write, reject  blogs/123/* write permission
	if s.isUserHasBroaderPolicy(ctx, policy) {
		return nil, ErrUserAlreadyHasBroaderPolicy
	}

	// Delete existing policies that match the resource prefix to avoid duplicates.
	// E.g., if adding blogs/* and user has blogs/123/*, remove blogs/123/* first.
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

func (s *service) isUserHasBroaderPolicy(ctx context.Context, policy *Policy) bool {
	currentPolicies, err := s.repo.Get(ctx, &GetPolicyRequest{
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

		if policy.Resource == resource {
			return true
		}

		if s.checkBroaderPolicy(policy.Resource, resource) {
			return true
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
		return strings.TrimSuffix(resource, "/*") // E.g., blogs/* -> blogs
	}

	if strings.HasSuffix(resource, "/") {
		return resource // E.g., blogs/
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
	if s.checkBroaderPolicy(policyResource, requestResource) {
		return true
	}

	// If user has access to specific resource and its sub-resources. E.g., blogs/11/pages/12
	// Then they can read blogs/11, but not write.
	if strings.HasPrefix(policyResource, requestResource) && action == ActionRead {
		return true
	}

	return false
}

func (s *service) checkBroaderPolicy(userResource, resourceRequested string) bool {
	if strings.HasSuffix(userResource, "/*") {
		prefix := strings.TrimSuffix(userResource, "*")     // E.g., blogs/* -> blogs/
		return strings.HasPrefix(resourceRequested, prefix) // E.g., blogs/123 has prefix blogs/
	}

	return false
}
