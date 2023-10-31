package rest

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"go.uber.org/dig"
	"microservice/app/core"
	"net/http"
	"os"
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

func InitImageServer(log core.Logger) {
	port := viper.GetString("img.port")
	imgPath := viper.GetString("img.path")

	// Also serve categories list
	entries, err := os.ReadDir(imgPath + "/categories")
	if err != nil {
		log.Fatal(err.Error())
	}

	categories := make(map[string][]string)
	for _, e := range entries {
		if e.IsDir() {
			if _, ok := categories[e.Name()]; !ok {
				categories[e.Name()] = []string{}
			}
		}
	}
	for c := range categories {
		entries, err := os.ReadDir(imgPath + "/categories/" + c)
		if err != nil {
			log.Fatal(err.Error())
		}
		for _, e := range entries {
			if !e.IsDir() {
				categories[c] = append(categories[c], e.Name())
			}
		}
	}

	fs := http.FileServer(http.Dir(imgPath))
	http.Handle("/static/", http.StripPrefix("/static", fs))
	http.HandleFunc("/static-list/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "*")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		json.NewEncoder(w).Encode(categories)
	})

	log.Info("Image server listening on port :" + port)
	err = http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal(err.Error())
	}
}
