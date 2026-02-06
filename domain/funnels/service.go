package funnels

import (
	"context"

	"github.com/adhikag24/policy-based-permission-model/domain/policies"
)

type Service interface {
	CreateFunnel(ctx context.Context, request *CreateFunnelRequest) error
	EditFunnel(ctx context.Context, request *EditFunnelRequest) error
	GetFunnel(ctx context.Context, request *GetFunnelRequest) (*Funnel, error)
}

type service struct {
	policiesService policies.Service
}

func NewService(policiesService policies.Service) Service {
	return &service{
		policiesService: policiesService,
	}
}

func (s *service) CreateFunnel(ctx context.Context, request *CreateFunnelRequest) error {
	if isPermitted := s.policiesService.CheckPermission(ctx, &policies.CheckPermissionRequest{
		AccountID:    request.AccountID,
		TeamMemberID: request.TeamMemberID,
		Resource:     "funnels/*",
		Action:       policies.ActionWrite,
	}); !isPermitted {
		return ErrPermissionDenied
	}

	return nil
}

func (s *service) EditFunnel(ctx context.Context, request *EditFunnelRequest) error {
	if isPermitted := s.policiesService.CheckPermission(ctx, &policies.CheckPermissionRequest{
		AccountID:    request.AccountID,
		TeamMemberID: request.TeamMemberID,
		Resource:     "funnels/" + request.FunnelID,
		Action:       policies.ActionWrite,
	}); !isPermitted {
		return ErrPermissionDenied
	}

	return nil
}

func (s *service) GetFunnel(ctx context.Context, request *GetFunnelRequest) (*Funnel, error) {
	if isPermitted := s.policiesService.CheckPermission(ctx, &policies.CheckPermissionRequest{
		AccountID:    request.AccountID,
		TeamMemberID: request.TeamMemberID,
		Resource:     "funnels/" + request.FunnelID,
		Action:       policies.ActionRead,
	}); !isPermitted {
		return nil, ErrPermissionDenied
	}

	return &Funnel{
		FunnelID: request.FunnelID,
		Name:     "Demo Funnel",
	}, nil
}
