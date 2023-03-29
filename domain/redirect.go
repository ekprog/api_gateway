package domain

import (
	"context"
	"microservice/app/core"
)

type RedirectUCase interface {
	Route(context.Context, *RedirectRouteRequest) (*RedirectRouteResponse, error)
}

//

type RedirectRouteRequest struct {
	AuthToken *string
	Address   string
	Data      []byte
}

type RedirectRouteResponse struct {
	Status   core.Status
	Response []byte
}
