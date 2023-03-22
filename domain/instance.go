package domain

import "context"

type Instance struct {
	Id       int64  `json:"id"`
	Name     string `json:"name"`
	Folder   string `json:"folder"`
	Endpoint string `json:"endpoint"`
	IsActive bool   `json:"is_active"`

	// after StatusService
	Status bool `json:"status"`
}

type InstancesRepository interface {
	All() ([]*Instance, error)
	GetByFolder(string) (*Instance, error)
	Insert(*Instance) error
	Delete(int64) error
	Update(instance *Instance) error
}

type InstancesInteractor interface {
	All(context.Context) (*InstancesAllResponse, error)
	Update(request UpdateInstanceRequest) (*StatusResponse, error)
}

// Delivery

type UpdateInstanceRequest struct {
	Id       *int64  `json:"id"`
	Folder   *string `json:"folder"`
	Endpoint *string `json:"endpoint"`
	IsActive *string `json:"is_active"`
}

type InstancesAllResponse struct {
	StatusCode string      `json:"status"`
	Instances  []*Instance `json:"services"`
}
