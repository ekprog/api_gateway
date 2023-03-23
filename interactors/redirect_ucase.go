package interactors

import (
	"context"
	"github.com/pkg/errors"
	"microservice/app/core"
	"microservice/domain"
	"microservice/services"
	"strconv"
)

type RedirectUCase struct {
	log           core.Logger
	routesRepo    domain.RoutesRepository
	authService   *services.AuthService
	callerService *services.ProtoCallerService
}

func NewRedirectUCase(log core.Logger,
	routesRepo domain.RoutesRepository,
	authService *services.AuthService,
	callerService *services.ProtoCallerService) *RedirectUCase {
	return &RedirectUCase{
		log:           log,
		routesRepo:    routesRepo,
		authService:   authService,
		callerService: callerService,
	}
}

// Route Перенаправляет входящий REST запрос на микросервис
func (ucase *RedirectUCase) Route(ctx context.Context, req *domain.RedirectRouteRequest) (*domain.RedirectRouteResponse, error) {
	if req == nil {
		return nil, errors.Errorf("empty request")
	}

	// Find same route in db
	route, err := ucase.routesRepo.GetByAddress(ctx, req.Address)
	if err != nil {
		return nil, errors.Wrapf(err, "error while fetching route for routing %s", req.Address)
	}
	if route == nil {
		return &domain.RedirectRouteResponse{
			Status: core.Status{
				Code: core.NotFound,
			},
		}, nil
	}

	// For call into Microservice
	callOptions := services.ProtoCall{
		Instance: route.Instance,
		Service:  route.ProtoService,
		Method:   route.ProtoMethod,
		Data:     req.Data,
		Headers:  map[string]string{},
	}

	// AUTH
	if route.AccessRole > core.RoleGuest {
		user, err := ucase.authService.Verify(ctx, req.AuthToken, route.AccessRole)
		if err != nil {
			return nil, errors.Wrapf(err, "error while verifying request")
		}
		if user == nil {
			return &domain.RedirectRouteResponse{
				Status: core.Status{
					Code: core.Unauthorised,
				},
			}, nil
		}
		// Set headers
		callOptions.Headers["user_id"] = strconv.FormatInt(user.Id, 10)
	}

	//
	// HERE REQUEST IS AUTHORIZED!
	//

	// Call
	response := &core.StatusResponse{}
	bytes, err := ucase.callerService.CallAndParse(ctx, callOptions, response)
	if err != nil {
		return nil, errors.Wrapf(err, "error while call instance method")
	}

	return &domain.RedirectRouteResponse{
		Response: bytes,
		Status: core.Status{
			Code: core.Success,
		},
	}, nil
}
