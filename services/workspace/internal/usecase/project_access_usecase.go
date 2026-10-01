package usecase

import (
	"context"
	"errors"

	"github.com/ricardocuellar/devboard/services/workspace/internal/domain"
)

type ProjectAccessResult struct {
	WorkspaceID string
}

type ProjectAccessValidator interface {
	Validate(
		ctx context.Context,
		projectID string,
		userID string,
	) (ProjectAccessResult, error)
}

type ProjectAccessProjectRepository interface {
	GetByID(
		ctx context.Context,
		id string,
	) (*domain.Project, error)
}

type ProjectAccessWorkspaceRepository interface {
	GetMember(
		ctx context.Context,
		workspaceID string,
		userID string,
	) (*domain.WorkspaceMember, error)
}

type ProjectAccessUseCase struct {
	projectRepository   ProjectAccessProjectRepository
	workspaceRepository ProjectAccessWorkspaceRepository
}

func NewProjectAccessUseCase(
	projectRepository ProjectAccessProjectRepository,
	workspaceRepository ProjectAccessWorkspaceRepository,
) *ProjectAccessUseCase {
	return &ProjectAccessUseCase{
		projectRepository: projectRepository,

		workspaceRepository: workspaceRepository,
	}
}

func (uc *ProjectAccessUseCase) Validate(
	ctx context.Context,
	projectID string,
	userID string,
) (ProjectAccessResult, error) {
	if projectID == "" ||
		userID == "" {
		return ProjectAccessResult{},
			domain.ErrInvalidInput
	}

	project, err :=
		uc.projectRepository.GetByID(
			ctx,
			projectID,
		)

	if err != nil {
		return ProjectAccessResult{},
			err
	}

	member, err :=
		uc.workspaceRepository.GetMember(
			ctx,
			project.WorkspaceID,
			userID,
		)

	if err != nil {
		if errors.Is(
			err,
			domain.ErrNotFound,
		) {
			return ProjectAccessResult{},
				domain.ErrForbidden
		}

		return ProjectAccessResult{},
			err
	}

	if member == nil {
		return ProjectAccessResult{},
			domain.ErrForbidden
	}

	return ProjectAccessResult{
		WorkspaceID: project.WorkspaceID,
	}, nil
}
