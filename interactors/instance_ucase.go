package interactors

import (
	"context"
	"github.com/pkg/errors"
	"microservice/app/core"
	"microservice/domain"
	"microservice/services"
	"microservice/tools"
)

type InstanceInteractor struct {
	log           core.Logger
	servicesRepo  domain.InstancesRepository
	statusService *services.StatusService
}

func NewInstanceInteractor(log core.Logger,
	repo domain.InstancesRepository,
	statusService *services.StatusService) *InstanceInteractor {
	return &InstanceInteractor{
		log:           log,
		servicesRepo:  repo,
		statusService: statusService,
	}
}

func (s *InstanceInteractor) All(ctx context.Context) (*domain.InstancesAllResponse, error) {
	items, err := s.servicesRepo.All(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "error while getting services list %s")
	}

	// ping all servers
	for _, item := range items {
		status, err := s.statusService.GetStatus(ctx, item.Folder)
		if err != nil {
			s.log.DebugWrap(err, "error while getting instance status %s", item.Folder)
			item.Status = false
		} else {
			item.Status = status
		}
	}

	return &domain.InstancesAllResponse{
		StatusCode: core.Success,
		Instances:  items,
	}, nil
}

func (s *InstanceInteractor) Update(ctx context.Context, req *tools.UpdateReq) (*core.StatusResponse, error) {

	err := s.servicesRepo.Update(ctx, req)
	if err != nil {
		return nil, errors.Wrap(err, "error while updating instance")
	}
	return &core.StatusResponse{
		Status: core.Status{
			Code: core.Success,
		},
	}, nil
}
