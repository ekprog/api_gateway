package jobs

import (
	"context"
	"microservice/app/core"
	"microservice/services"
	"time"
)

type GetServicesStatusesJob struct {
	log           core.Logger
	statusService *services.StatusService
}

func NewGetServicesStatusesJob(
	log core.Logger,
	statusService *services.StatusService) *GetServicesStatusesJob {
	return &GetServicesStatusesJob{
		log:           log,
		statusService: statusService,
	}
}

func (j *GetServicesStatusesJob) Run() error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	err := j.statusService.CheckAll(ctx)
	if err != nil {
		j.log.ErrorWrap(err, "error in job GetServicesStatusesJob")
	} else {
		j.log.Info("GetServicesStatusesJob successfully finished work!")
	}

	return nil
}
