package bootstrap

import (
	"github.com/pkg/errors"
	"go.uber.org/dig"
	"microservice/app/rest"
	"microservice/delivery"
	"microservice/domain"
	"microservice/interactors"
	"microservice/repos"
	"microservice/services"
)

func initDependencies(di *dig.Container) error {

	// Repositories
	_ = di.Provide(
		repos.NewInstancesRepo,
		dig.As(new(domain.InstancesRepository)),
	)

	_ = di.Provide(
		repos.NewRoutesRepo,
		dig.As(new(domain.RoutesRepository)),
	)

	// Services
	_ = di.Provide(services.NewEndpointConnectionService)
	_ = di.Provide(services.NewAuthService)
	_ = di.Provide(services.NewStatusService)
	_ = di.Provide(services.NewProtoCallerService)

	// Use cases
	_ = di.Provide(
		interactors.NewInstanceInteractor,
		dig.As(new(domain.InstancesUCase)),
	)

	_ = di.Provide(
		interactors.NewRedirectUCase,
		dig.As(new(domain.RedirectUCase)),
	)

	// Delivery init
	err := rest.InitDelivery("/admin", delivery.NewAdminDelivery)
	if err != nil {
		return err
	}

	if err := rest.InitDeliveryDynamic("/api/v1", delivery.NewRouterDelivery); err != nil {
		return errors.Wrap(err, "error while route gateway router")
	}

	return nil
}
