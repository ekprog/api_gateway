package delivery

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"microservice/app/core"
	"microservice/domain"
	"microservice/services"
	"strconv"
)

type AdminDelivery struct {
	log            core.Logger
	instancesUCase domain.InstancesInteractor
	routesRepo     domain.RoutesRepository

	authService         *services.AuthService
	endpointConnService *services.EndpointConnectionService
}

func NewAdminDelivery(log core.Logger,
	instancesUCase domain.InstancesInteractor,
	authService *services.AuthService) *AdminDelivery {
	return &AdminDelivery{
		log:            log,
		instancesUCase: instancesUCase,
		authService:    authService,
	}
}

func (d *AdminDelivery) Route(g *gin.RouterGroup) error {
	g.Use(d.ErrorMW)
	g.Use(d.GeneralMW)
	g.Use(d.AdminAuthMW)

	g.POST("/services", d.Services)
	g.POST("/services/update", d.Update)

	return nil
}

func (d *AdminDelivery) GeneralMW(c *gin.Context) {
	c.Header("content-type", "application/json")
}

func (d *AdminDelivery) ErrorMW(c *gin.Context) {
	c.Next()
	if len(c.Errors) > 0 {
		err := c.Errors[0].Err
		d.log.ErrorWrap(err, "Request error")
		c.JSON(500, domain.ServerError())
	}
}

func (d *AdminDelivery) AdminAuthMW(c *gin.Context) {
	user, err := d.authService.VerifyRequest(c, domain.RoleSuperAdmin)
	if err != nil {
		_ = c.Error(errors.Wrap(err, "error while verifying admin request"))
		c.JSON(500, domain.ServerError())
		return
	}
	if user == nil {
		c.JSON(500, domain.UnauthorizedError())
		return
	}
	c.Header("user_id", strconv.FormatInt(user.Id, 10))
	c.Next()
}

func (d *AdminDelivery) Services(c *gin.Context) {
	res, err := d.instancesUCase.All(c)
	if err != nil {
		_ = c.AbortWithError(500, errors.Wrap(err, "error while services_all ucase"))
		return
	}
	c.JSON(200, res)
}

func (d *AdminDelivery) Update(c *gin.Context) {

	request := domain.UpdateInstanceRequest{}

	err := json.NewDecoder(c.Request.Body).Decode(&request)
	if err != nil {
		_ = c.AbortWithError(500, errors.Wrap(err, "validation error"))
		c.JSON(500, domain.ValidationError())
		return
	}

	res, err := d.instancesUCase.Update(request)
	if err != nil {
		_ = c.AbortWithError(500, errors.Wrap(err, "error while services_all ucase"))
		return
	}
	c.JSON(200, res)
}
