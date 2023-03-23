package domain

import (
	"context"
	"microservice/app/core"
)

type Route struct {
	Id           int64
	HttpMethod   string
	HttpAddress  string
	Instance     string
	ProtoService string
	ProtoMethod  string
	AccessRole   core.AccessRole
	IsActive     bool
}

type RoutesRepository interface {
	All(context.Context) ([]*Route, error)
	GetByAddress(ctx context.Context, addr string) (*Route, error)
	Insert(context.Context, *Route) error
	Delete(context.Context, int64) error
}
