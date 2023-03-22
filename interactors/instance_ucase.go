package interactors

import (
	"context"
	"github.com/pkg/errors"
	"microservice/app/core"
	"microservice/domain"
	"microservice/services"
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
	items, err := s.servicesRepo.All()
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
		StatusCode: domain.Success,
		Instances:  items,
	}, nil
}

func (s *InstanceInteractor) Update(request domain.UpdateInstanceRequest) (*domain.StatusResponse, error) {
	//err := s.servicesRepo.Update(instance)
	//if err != nil {
	//	return nil, errors.Wrap(err, "error while updating instance")
	//}
	return &domain.StatusResponse{
		Status: domain.Status{
			Code:    domain.Success,
			Message: domain.Success,
		},
	}, nil
}
