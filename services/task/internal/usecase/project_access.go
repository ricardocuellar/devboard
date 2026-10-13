package usecase

import "context"

type ProjectAccess struct {
	WorkpaceID string
}

type ProjectAccessChecker interface {
	ValidateProjectAccess(ctx context.Context, projectID string, userID string) (ProjectAccess, error)
}
