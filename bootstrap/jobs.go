package bootstrap

import (
	"microservice/app/job"
	"microservice/jobs"
)

func initJobs() {
	job.NewJob(jobs.NewGetServicesStatusesJob, job.Time("5 minutes"))
}
