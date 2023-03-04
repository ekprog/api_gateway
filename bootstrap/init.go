package bootstrap

import (
	"api_gateway/app"
	"api_gateway/delivery"
	"api_gateway/services"
	"database/sql"
)

func injectDependencies(db *sql.DB, logger app.Logger) error {

	// DI
	routerService := services.NewRouterService(logger)

	// Delivery Init
	routeDelivery := delivery.NewRouterDelivery(logger, routerService)

	// Register
	err := app.InitRestDelivery("/api/v1", routeDelivery)
	if err != nil {
		return err
	}
	return nil
}

//func provide(diObj *dig.Container, list ...interface{}) error {
//	for _, p := range list {
//		if err := diObj.Provide(p); err != nil {
//			return err
//		}
//	}
//	return nil
//}
