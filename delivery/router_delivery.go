package delivery

import (
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"microservice/app/core"
	"microservice/domain"
	"microservice/services"
	"strconv"
)

type RouterDelivery struct {
	log           core.Logger
	routesRepo    domain.RoutesRepository
	authService   *services.AuthService
	callerService *services.ProtoCallerService
}

func NewRouterDelivery(log core.Logger,
	routesRepo domain.RoutesRepository,
	authService *services.AuthService,
	callerService *services.ProtoCallerService,
) *RouterDelivery {
	return &RouterDelivery{
		log:           log,
		routesRepo:    routesRepo,
		authService:   authService,
		callerService: callerService,
	}
}

func (d *RouterDelivery) Route(c *gin.Context) {

	c.Header("content-type", "application/json")

	// Example: /auth/login
	url := c.Request.URL.Path

	// Find same route in db
	route, err := d.routesRepo.GetByAddress(url)
	if err != nil {
		d.log.ErrorWrap(err, "error while fetching route for routing %s", url)
		c.JSON(500, domain.ServerError())
		return
	}
	if route == nil {
		d.log.Error("cannot find route for routing %s", url)
		c.JSON(500, domain.NotFoundError())
		return
	}

	// For call into Microservice
	callOptions := services.ProtoCall{
		Instance: route.Instance,
		Service:  route.ProtoService,
		Method:   route.ProtoMethod,
		Data:     nil,
		Headers:  map[string]string{},
	}

	// AUTH
	if route.AccessRole > domain.RoleGuest {
		user, err := d.authService.VerifyRequest(c, route.AccessRole)
		if err != nil {
			_ = c.Error(errors.Wrap(err, "error while verifying request"))
			c.JSON(500, domain.ServerError())
			return
		}
		if user == nil {
			c.JSON(500, domain.UnauthorizedError())
			return
		}
		// Set headers
		callOptions.Headers["user_id"] = strconv.FormatInt(user.Id, 10)
	}

	//
	// HERE REQUEST IS AUTHORIZED!
	//

	// Call
	response := &domain.StatusResponse{}
	bytes, err := d.callerService.CallAndParse(c, callOptions, response)
	if err != nil {
		_ = c.Error(errors.Wrapf(err, "in instance error (%s)", route.Instance))
		c.JSON(500, domain.ServerError())
		return
	}

	// Status
	if response == nil || response.Status.Code != domain.Success {
		c.Status(500)
	} else {
		c.Status(200)
	}

	// To client
	c.Writer.Write(bytes)
}
