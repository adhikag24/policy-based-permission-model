package blogs

import (
	"context"
	"fmt"

	"github.com/adhikag24/policy-based-permission-model/domain/policies"
)

type Service interface {
	WriteBlogPage(ctx context.Context, request *WriteBlogPageRequest) error
	ReadBlogPage(ctx context.Context, request *ReadBlogPageRequest) error
	ReadBlogSettings(ctx context.Context, request *ReadBlogSettingsRequest) error
	WriteBlogSettings(ctx context.Context, request *WriteBlogSettingsRequest) error
}

type service struct {
	policiesService policies.Service
}

func NewService(policiesService policies.Service) Service {
	return &service{
		policiesService: policiesService,
	}
}

func (s *service) ReadBlogSettings(ctx context.Context, request *ReadBlogSettingsRequest) error {
	if isPermitted := s.policiesService.CheckPermission(ctx, &policies.CheckPermissionRequest{
		AccountID:    request.AccountID,
		TeamMemberID: request.TeamMemberID,
		Resource:     fmt.Sprintf("blogs/%s/settings", request.BlogID),
		Action:       policies.ActionRead,
	}); !isPermitted {
		return ErrPermissionDenied
	}

	return nil
}

func (s *service) WriteBlogPage(ctx context.Context, request *WriteBlogPageRequest) error {
	if isPermitted := s.policiesService.CheckPermission(ctx, &policies.CheckPermissionRequest{
		AccountID:    request.AccountID,
		TeamMemberID: request.TeamMemberID,
		Resource:     fmt.Sprintf("blogs/%s", request.PageID), // Check if user has permission to write this blog page.
		Action:       policies.ActionWrite,
	}); !isPermitted {
		return ErrPermissionDenied
	}

	return nil
}

func (s *service) WriteBlogSettings(ctx context.Context, request *WriteBlogSettingsRequest) error {
	if isPermitted := s.policiesService.CheckPermission(ctx, &policies.CheckPermissionRequest{
		AccountID:    request.AccountID,
		TeamMemberID: request.TeamMemberID,
		Resource:     fmt.Sprintf("blogs/%s/settings", request.BlogID),
		Action:       policies.ActionWrite,
	}); !isPermitted {
		return ErrPermissionDenied
	}

	return nil
}

func (s *service) ReadBlogPage(ctx context.Context, request *ReadBlogPageRequest) error {
	if isPermitted := s.policiesService.CheckPermission(ctx, &policies.CheckPermissionRequest{
		AccountID:    request.AccountID,
		TeamMemberID: request.TeamMemberID,
		Resource:     fmt.Sprintf("blogs/%s/pages/%s", request.BlogID, request.PageID), // Simulates multiple identifiers in resource.
		Action:       policies.ActionRead,
	}); !isPermitted {
		return ErrPermissionDenied
	}

	return nil
}
