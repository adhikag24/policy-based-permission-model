package policies_test

import (
	"testing"

	"github.com/adhikag24/policy-based-permission-model/domain/policies"
	mockRepository "github.com/adhikag24/policy-based-permission-model/domain/policies/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type test struct {
	mockRepository *mockRepository.MockRepository
}

func setup(ctrl *gomock.Controller) *test {
	return &test{
		mockRepository: mockRepository.NewMockRepository(ctrl),
	}
}

func TestCreatePolicy(t *testing.T) {
	t.Run("Successfully creates policy when user already has broader policy", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		test := setup(ctrl)
		test.mockRepository.EXPECT().Get(gomock.Any(), &policies.GetPolicyRequest{
			AccountID:    100,
			TeamMemberID: 200,
			Action:       policies.ActionRead,
		}).Return([]policies.Policy{
			{
				ID:           1,
				AccountID:    100,
				TeamMemberID: 200,
				Resource:     "funnels/123/pages/123/*",
				Action:       policies.ActionRead,
			},
		}, nil)
		service := policies.NewService(test.mockRepository)

		policy, err := service.CreatePolicy(t.Context(), &policies.Policy{
			AccountID:    100,
			TeamMemberID: 200,
			Resource:     "funnels/123/pages/123/components/456",
			Action:       policies.ActionRead,
		})

		assert.ErrorIs(t, err, policies.ErrUserAlreadyHasBroaderPolicy)
		assert.Nil(t, policy)
	})

	t.Run("Successfully creates policy when no broader policy exists and overrides lower policy", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		test := setup(ctrl)
		test.mockRepository.EXPECT().Get(gomock.Any(), &policies.GetPolicyRequest{
			AccountID:    100,
			TeamMemberID: 200,
			Action:       policies.ActionRead,
		}).Return([]policies.Policy{
			{
				ID:           1,
				AccountID:    100,
				TeamMemberID: 200,
				Resource:     "articles/*",
				Action:       policies.ActionRead,
			},
			{
				ID:           1,
				AccountID:    100,
				TeamMemberID: 200,
				Resource:     "blogs/123/*",
				Action:       policies.ActionRead,
			},
		}, nil)
		test.mockRepository.EXPECT().DeleteByPrefix(gomock.Any(), &policies.DeleteByPrefixRequest{
			AccountID:      100,
			TeamMemberID:   200,
			ResourcePrefix: "blogs/",
			Action:         policies.ActionRead,
		}).Return(nil)
		test.mockRepository.EXPECT().Create(gomock.Any(), &policies.Policy{
			AccountID:    100,
			TeamMemberID: 200,
			Resource:     "blogs/*",
			Action:       policies.ActionRead,
		}).Return(&policies.Policy{
			ID:           1,
			AccountID:    100,
			TeamMemberID: 200,
			Resource:     "blogs/*",
			Action:       policies.ActionRead,
		}, nil)
		service := policies.NewService(test.mockRepository)

		policy, err := service.CreatePolicy(t.Context(), &policies.Policy{
			AccountID:    100,
			TeamMemberID: 200,
			Resource:     "blogs/*",
			Action:       policies.ActionRead,
		})

		assert.NoError(t, err)
		assert.NotNil(t, policy)
		assert.Equal(t, int64(1), policy.ID)
	})
}

func TestCheckPermission(t *testing.T) {
	tests := []struct {
		name          string
		request       *policies.CheckPermissionRequest
		mockPolicies  []policies.Policy
		mockAction    policies.Action
		wantPermitted bool
	}{
		{
			name: "permission granted for exact match",
			request: &policies.CheckPermissionRequest{
				AccountID:    100,
				TeamMemberID: 200,
				Resource:     "projects/123/tasks/456",
				Action:       policies.ActionRead,
			},
			mockAction: policies.ActionRead,
			mockPolicies: []policies.Policy{
				{
					ID:           1,
					AccountID:    100,
					TeamMemberID: 200,
					Resource:     "projects/123/tasks/456",
					Action:       policies.ActionRead,
				},
			},
			wantPermitted: true,
		},
		{
			name: "permission granted when broader policy exists",
			request: &policies.CheckPermissionRequest{
				AccountID:    100,
				TeamMemberID: 200,
				Resource:     "projects/123/tasks/456",
				Action:       policies.ActionRead,
			},
			mockAction: policies.ActionRead,
			mockPolicies: []policies.Policy{
				{
					ID:           1,
					AccountID:    100,
					TeamMemberID: 200,
					Resource:     "projects/*",
					Action:       policies.ActionRead,
				},
			},
			wantPermitted: true,
		},
		{
			name: "permission denied when no matching policy exists",
			request: &policies.CheckPermissionRequest{
				AccountID:    100,
				TeamMemberID: 200,
				Resource:     "projects/123/tasks/456",
				Action:       policies.ActionWrite,
			},
			mockAction: policies.ActionWrite,
			mockPolicies: []policies.Policy{
				{
					ID:           1,
					AccountID:    100,
					TeamMemberID: 200,
					Resource:     "projects/123",
					Action:       policies.ActionRead,
				},
			},
			wantPermitted: false,
		},
		{
			name: "permission granted with root policy",
			request: &policies.CheckPermissionRequest{
				AccountID:    100,
				TeamMemberID: 200,
				Resource:     "projects/999",
				Action:       policies.ActionRead,
			},
			mockAction: policies.ActionRead,
			mockPolicies: []policies.Policy{
				{
					ID:           1,
					AccountID:    100,
					TeamMemberID: 200,
					Resource:     "*",
					Action:       policies.ActionRead,
				},
			},
			wantPermitted: true,
		},
		{
			name: "read permitted on parent resource when child policy exists",
			request: &policies.CheckPermissionRequest{
				AccountID:    100,
				TeamMemberID: 200,
				Resource:     "projects/123",
				Action:       policies.ActionRead,
			},
			mockAction: policies.ActionRead,
			mockPolicies: []policies.Policy{
				{
					ID:           1,
					AccountID:    100,
					TeamMemberID: 200,
					Resource:     "projects/123/tasks/456",
					Action:       policies.ActionRead,
				},
			},
			wantPermitted: true,
		},
		{
			name: "write denied on parent resource when child policy exists",
			request: &policies.CheckPermissionRequest{
				AccountID:    100,
				TeamMemberID: 200,
				Resource:     "projects/123",
				Action:       policies.ActionWrite,
			},
			mockAction: policies.ActionWrite,
			mockPolicies: []policies.Policy{
				{
					ID:           1,
					AccountID:    100,
					TeamMemberID: 200,
					Resource:     "projects/123/tasks/456",
					Action:       policies.ActionWrite,
				},
			},
			wantPermitted: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			test := setup(ctrl)
			test.mockRepository.EXPECT().Get(gomock.Any(), &policies.GetPolicyRequest{
				AccountID:    tt.request.AccountID,
				TeamMemberID: tt.request.TeamMemberID,
				Action:       tt.mockAction,
			}).Return(tt.mockPolicies, nil)

			service := policies.NewService(test.mockRepository)

			hasPermission := service.CheckPermission(t.Context(), tt.request)
			assert.Equal(t, tt.wantPermitted, hasPermission)
		})
	}
}
