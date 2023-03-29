package delivery

import (
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"microservice/app/core"
	"microservice/app/rest"
	"microservice/delivery/forms"
	"microservice/domain"
	"microservice/services"
	"microservice/tools"
	"strconv"
)

type AdminDelivery struct {
	log            core.Logger
	instancesUCase domain.InstancesUCase
	routesRepo     domain.RoutesRepository

	authService         *services.AuthService
	endpointConnService *services.EndpointConnectionService
}

func NewAdminDelivery(log core.Logger,
	instancesUCase domain.InstancesUCase,
	authService *services.AuthService) *AdminDelivery {
	return &AdminDelivery{
		log:            log,
		instancesUCase: instancesUCase,
		authService:    authService,
	}
}

func (d *AdminDelivery) Route(g *gin.RouterGroup) error {
	g.Use(rest.GeneralMW)
	g.Use(rest.ErrorMW)
	g.Use(d.AdminAuthMW)

	g.POST("/services", d.Services)
	g.POST("/services/update", d.Update)

	return nil
}

func (d *AdminDelivery) AdminAuthMW(ctx *gin.Context) {

	// Extract token
	authTokens := ctx.Request.Header["Authorization"]
	if authTokens == nil || len(authTokens) <= 0 {
		ctx.AbortWithStatusJSON(500, rest.UnauthorizedError())
		return
	}
	authToken := authTokens[0]
	d.log.Debug("Authorization access with token: %v", authToken)

	user, err := d.authService.Verify(ctx, authToken, core.RoleSuperAdmin)
	if err != nil {
		_ = ctx.Error(errors.Wrap(err, "error while verifying admin request"))
		ctx.AbortWithStatusJSON(500, rest.ServerError())
		return
	}
	if user == nil {
		ctx.AbortWithStatusJSON(500, rest.UnauthorizedError())
		return
	}
	ctx.Header("user_id", strconv.FormatInt(int64(user.Id), 10))
	ctx.Next()
}

func (d *AdminDelivery) Services(ctx *gin.Context) {
	res, err := d.instancesUCase.All(ctx)
	if err != nil {
		_ = ctx.AbortWithError(500, errors.Wrap(err, "error while services_all ucase"))
		return
	}
	ctx.JSON(200, res)
}

func (d *AdminDelivery) Update(ctx *gin.Context) {

	// Validation
	reqObj := &forms.InstanceUpdateForm{}
	err := ctx.BindJSON(reqObj)
	if err != nil {
		ctx.Abort() // err already in context
		return
	}

	// To Update request
	updateReq, err := tools.NewUpdateReqReader(ctx.Request.Body)
	if err != nil {
		ctx.Abort()
		return
	}

	//
	res, err := d.instancesUCase.Update(ctx, updateReq)
	if err != nil {
		_ = ctx.AbortWithError(500, errors.Wrap(err, "error while services_all ucase"))
		return
	}

	//
	ctx.JSON(200, res)
}
