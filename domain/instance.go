package domain

import (
	"context"
	"microservice/app/core"
	"microservice/tools"
)

type Instance struct {
	Id       int32  `json:"id"`
	Name     string `json:"name"`
	Folder   string `json:"folder"`
	Endpoint string `json:"endpoint"`
	IsActive bool   `json:"is_active"`
	Status   bool   `json:"status"`
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
