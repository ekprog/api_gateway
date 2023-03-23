package domain

import (
	"context"
	"microservice/app/core"
	"microservice/tools"
)

type Instance struct {
	Id       int32
	Name     string
	Folder   string
	Endpoint string
	IsActive bool

	// after StatusService
	Status bool
}

type InstancesRepository interface {
	All(context.Context) ([]*Instance, error)
	GetByFolder(context.Context, string) (*Instance, error)
	Insert(context.Context, *Instance) error
	Delete(context.Context, int32) error
	Update(context.Context, *tools.UpdateReq) error
}

type InstancesUCase interface {
	All(context.Context) (*InstancesAllResponse, error)
	Update(context.Context, *tools.UpdateReq) (*core.StatusResponse, error)
}

// Delivery
type InstancesAllResponse struct {
	StatusCode string      `json:"status"`
	Instances  []*Instance `json:"services"`
}
