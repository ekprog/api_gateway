package delivery

import (
	"api_gateway/app"
	"api_gateway/domain"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

type RouterDelivery struct {
	log           app.Logger
	routerService domain.RouterService
}

func NewRouterDelivery(log app.Logger, routerService domain.RouterService) *RouterDelivery {
	return &RouterDelivery{
		log:           log,
		routerService: routerService,
	}
}

func (d *RouterDelivery) Route(g *gin.RouterGroup) error {

	d.routerService.SetRouter(g)

	// Init Services
	err := d.routerService.MakeService(domain.Service{
		ProtoServiceName: "AuthService",
		ProtoDir:         "./proto/auth_service",
		ProtoFilename:    "api/auth_service.proto",
		HttpAddress:      "localhost:8086",
	})
	if err != nil {
		return errors.Wrap(err, "cannot make Test service")
	}

	err = d.routerService.MakeService(domain.Service{
		ProtoServiceName: "TestService",
		ProtoDir:         "./proto/test_service",
		ProtoFilename:    "api/test_service.proto",
		HttpAddress:      "localhost:8087",
	})
	if err != nil {
		return errors.Wrap(err, "cannot make Test service")
	}

	err = d.routerService.MakeService(domain.Service{
		ProtoServiceName: "ToDoService",
		ProtoDir:         "./proto/todo_service",
		ProtoFilename:    "api/service.proto",
		HttpAddress:      "localhost:8071",
	})
	if err != nil {
		return errors.Wrap(err, "cannot make ToDo service")
	}

	err = d.routerService.MakeService(domain.Service{
		ProtoServiceName: "ProfilesService",
		ProtoDir:         "./proto/profiles_service",
		ProtoFilename:    "api/service.proto",
		HttpAddress:      "localhost:8072",
	})
	if err != nil {
		return errors.Wrap(err, "cannot make ToDo service")
	}

	// Auth
	err = d.routerService.Handle(domain.Route{
		HttpMethod:   "POST",
		HttpAddress:  "/auth/register",
		ProtoService: "AuthService",
		ProtoMethod:  "Register",
		AccessRole:   domain.RoleGuest,
	})
	if err != nil {
		return errors.Wrap(err, "cannot register route")
	}

	err = d.routerService.Handle(domain.Route{
		HttpMethod:   "POST",
		HttpAddress:  "/auth/login",
		ProtoService: "AuthService",
		ProtoMethod:  "Login",
		AccessRole:   domain.RoleGuest,
	})
	if err != nil {
		return errors.Wrap(err, "cannot register route")
	}

	err = d.routerService.Handle(domain.Route{
		HttpMethod:   "POST",
		HttpAddress:  "/auth/revoke",
		ProtoService: "AuthService",
		ProtoMethod:  "Revoke",
		AccessRole:   domain.RoleGuest,
	})
	if err != nil {
		return errors.Wrap(err, "cannot register route")
	}

	err = d.routerService.Handle(domain.Route{
		HttpMethod:   "POST",
		HttpAddress:  "/auth/verify",
		ProtoService: "AuthService",
		ProtoMethod:  "Verify",
		AccessRole:   domain.RoleGuest,
	})
	if err != nil {
		return errors.Wrap(err, "cannot register route")
	}

	err = d.routerService.Handle(domain.Route{
		HttpMethod:   "POST",
		HttpAddress:  "/auth/refresh",
		ProtoService: "AuthService",
		ProtoMethod:  "Refresh",
		AccessRole:   domain.RoleGuest,
	})
	if err != nil {
		return errors.Wrap(err, "cannot register route")
	}

	err = d.routerService.Handle(domain.Route{
		HttpMethod:   "POST",
		HttpAddress:  "/auth/list",
		ProtoService: "AuthService",
		ProtoMethod:  "List",
		AccessRole:   domain.RoleGuest,
	})
	if err != nil {
		return errors.Wrap(err, "cannot register route")
	}

	// Set Auth service
	err = d.routerService.SetAuthService("AuthService")
	if err != nil {
		return errors.Wrap(err, "cannot add auth service")
	}

	// ROUTING
	err = d.routerService.Handle(domain.Route{
		HttpMethod:   "POST",
		HttpAddress:  "/test/private",
		ProtoService: "TestService",
		ProtoMethod:  "TestPrivate",
		AccessRole:   domain.RoleUser,
	})
	if err != nil {
		return errors.Wrap(err, "cannot register route")
	}

	err = d.routerService.Handle(domain.Route{
		HttpMethod:   "POST",
		HttpAddress:  "/test/public",
		ProtoService: "TestService",
		ProtoMethod:  "TestPublic",
		AccessRole:   domain.RoleGuest,
	})
	if err != nil {
		return errors.Wrap(err, "cannot register route")
	}

	// TODOService
	err = d.routerService.Handle(domain.Route{
		HttpMethod:   "POST",
		HttpAddress:  "/projects/create",
		ProtoService: "ToDoService",
		ProtoMethod:  "CreateProject",
		AccessRole:   domain.RoleUser,
	})
	if err != nil {
		return errors.Wrap(err, "cannot register route")
	}

	// ProfilesService
	err = d.routerService.Handle(domain.Route{
		HttpMethod:   "POST",
		HttpAddress:  "/profile",
		ProtoService: "ProfilesService",
		ProtoMethod:  "Get",
		AccessRole:   domain.RoleGuest,
	})
	if err != nil {
		return errors.Wrap(err, "cannot register route")
	}

	err = d.routerService.Handle(domain.Route{
		HttpMethod:   "POST",
		HttpAddress:  "/profiles/update",
		ProtoService: "ProfilesService",
		ProtoMethod:  "Update",
		AccessRole:   domain.RoleGuest,
	})
	if err != nil {
		return errors.Wrap(err, "cannot register route")
	}

	return nil
}
