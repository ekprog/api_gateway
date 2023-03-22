package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"go.uber.org/dig"
	"microservice/app/core"
	"reflect"
	"runtime"
	"strings"
)

type DeliveryService interface {
	Route(r *gin.RouterGroup) error
}

type DeliveryDynamicService interface {
	Route(ctx *gin.Context)
}

func InitDelivery(path string, f interface{}) error {

	name := runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()

	di := core.GetDI()
	scope := di.Scope(name)

	err := scope.Provide(f, dig.As(new(DeliveryService)))
	if err != nil {
		return errors.Wrap(err, "cannot init rest delivery")
	}

	return scope.Invoke(func(d DeliveryService) error {
		g := restServer.Group(path)
		err := d.Route(g)
		if err != nil {
			return errors.Wrap(err, "error in route rest delivery")
		}
		return nil
	})
}

func InitDeliveryDynamic(path string, f interface{}) error {

	name := runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()

	di := core.GetDI()
	scope := di.Scope(name)

	err := scope.Provide(f, dig.As(new(DeliveryDynamicService)))
	if err != nil {
		return errors.Wrap(err, "cannot init rest delivery")
	}

	return scope.Invoke(func(d DeliveryDynamicService) error {
		restServer.NoRoute(func(ctx *gin.Context) {
			url := ctx.Request.URL.Path
			if strings.HasPrefix(url, path) {
				url = strings.TrimPrefix(url, path)
				url = strings.TrimSuffix(url, "/")
				ctx.Request.URL.Path = url
				ctx.Status(200) // replacing 404 with 200
				d.Route(ctx)
			}
		})
		return nil
	})

}
