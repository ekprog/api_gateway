package services

import (
	"context"
	"github.com/pkg/errors"
	"microservice/app/core"
	"microservice/domain"
)

// StatusService проверяет статус микросервисов, отправляя им ping периодически
type StatusService struct {
	log           core.Logger
	instancesRepo domain.InstancesRepository
	callerService *ProtoCallerService

	cache  map[string]bool
	errors map[string]string
}

func NewStatusService(log core.Logger,
	instancesRepo domain.InstancesRepository,
	callerService *ProtoCallerService) *StatusService {
	return &StatusService{
		log:           log,
		instancesRepo: instancesRepo,
		callerService: callerService,

		cache: map[string]bool{},

		// Причины недоступности
		errors: map[string]string{},
	}
}

// CheckAll checks every service for itself status
func (s *StatusService) CheckAll(ctx context.Context) error {
	items, err := s.instancesRepo.All(ctx)
	if err != nil {
		return errors.Wrap(err, "cannot check services for status")
	}
	for _, item := range items {
		status, err := s.GetStatus(ctx, item.Name)
		if err != nil {
			s.cache[item.Folder] = status
			s.errors[item.Folder] = err.Error()
		} else {
			s.cache[item.Folder] = status // only read here (no need for mutex)
			s.errors[item.Folder] = ""
		}
	}
	return nil
}

func (s *StatusService) GetStatus(ctx context.Context, instanceName string) (bool, error) {

	if cache, ok := s.cache[instanceName]; ok {
		return cache, nil
	}

	callOptions := ProtoCall{
		Instance: instanceName,
		Service:  "StatusService",
		Method:   "Ping",
	}

	response := &core.StatusResponse{}
	_, err := s.callerService.CallAndParse(ctx, callOptions, response)
	if err != nil {
		return false, errors.Wrapf(err, "cannot check status of %s", instanceName)
	}
	if response == nil || response.Status.Code != core.Success {
		return false, errors.Errorf("not success status code")
	}
	return true, nil
}
