package delivery

import (
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"io"
	"microservice/app/core"
	"microservice/domain"
)

type RouterDelivery struct {
	log         core.Logger
	routerUCase domain.RedirectUCase
}

func NewRouterDelivery(log core.Logger,
	routerUCase domain.RedirectUCase,
) *RouterDelivery {
	return &RouterDelivery{
		log:         log,
		routerUCase: routerUCase,
	}
}

func (d *RouterDelivery) Route(ctx *gin.Context) {

	// Headers
	ctx.Header("content-type", "application/json")

	// Extract token
	authTokens := ctx.Request.Header["Authorization"]
	if authTokens == nil || len(authTokens) <= 0 {
		ctx.JSON(500, core.StatusResponse{
			Status: core.Status{
				Code: core.Unauthorised,
			}})
		return
	}
	authToken := authTokens[0]
	d.log.Debug("Authorization access with token: %s", authToken)

	// Parse body
	body, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		_ = ctx.Error(errors.Wrapf(err, "cannot parse request`s body"))
		ctx.AbortWithStatus(500)
		return
	}

	// UCase
	res, err := d.routerUCase.Route(ctx, &domain.RedirectRouteRequest{
		AuthToken: authToken,
		Address:   ctx.Request.URL.Path,
		Data:      body,
	})
	if err != nil {
		_ = ctx.Error(errors.Wrapf(err, "cannot route client`s request"))
		ctx.AbortWithStatus(500)
		return
	}

	// Status
	if res == nil || res.Status.Code != core.Success || res.Response == nil {
		ctx.Status(500)
	} else {
		ctx.Status(200)
	}

	// To client
	ctx.Writer.Write(res.Response)
}
