package bootstrap

import (
	"database/sql"
	"github.com/pkg/errors"
	"microservice/app"
	"microservice/app/core"
	"microservice/app/job"
	"microservice/app/kafka"
	"microservice/app/rest"
)

func Run(rootPath ...string) error {

	// ENV, etc
	ctx, _, err := app.InitApp(rootPath...)
	if err != nil {
		return errors.Wrap(err, "error while init app")
	}

	// Logger
	logger, err := app.InitLogs(rootPath...)
	if err != nil {
		return errors.Wrap(err, "error while init logs")
	}

	// Storage
	err = app.InitStorage()
	if err != nil {
		return errors.Wrap(err, "error while init storage")
	}

	// Database
	db, err := app.InitDatabase()
	if err != nil {
		return errors.Wrap(err, "error while init db")
	}

	// Migrations
	err = app.RunMigrations(rootPath...)
	if err != nil {
		return errors.Wrap(err, "error while making migrations")
	}

	// gRPC
	err = rest.Init()
	if err != nil {
		return errors.Wrap(err, "cannot init gRPC")
	}

	// DI
	di := core.GetDI()

	if err = di.Provide(func() *sql.DB {
		return db
	}); err != nil {
		return errors.Wrap(err, "cannot provide db")
	}

	if err = di.Provide(func() core.Logger {
		return logger
	}); err != nil {
		return errors.Wrap(err, "cannot provide logger")
	}

	// CRON
	err = job.Init(logger, di)
	if err != nil {
		return errors.Wrap(err, "cannot init jobs")
	}

	// KAFKA
	err = kafka.InitKafka(logger)
	if err != nil {
		return errors.Wrap(err, "cannot init kafka")
	}

	// PROTO REGISTRY
	registry := app.NewProtoRegistry()
	err = registry.Init()
	if err != nil {
		return errors.Wrap(err, "cannot initialize proto registry")
	}
	if err = di.Provide(func() *app.ProtoRegistry {
		return registry
	}); err != nil {
		return errors.Wrap(err, "cannot provide proto registry")
	}

	// CORE
	if err := initDependencies(di); err != nil {
		return errors.Wrap(err, "error while init dependencies")
	}

	//
	//
	// HERE CORE READY FOR WORK...
	//
	//

	// CRON
	initJobs()

	if err := job.Start(); err != nil {
		return errors.Wrap(err, "error while start jobs")
	}

	// Run gRPC and block
	go rest.RunServer()
	go rest.InitImageServer(logger)

	// End context
	<-ctx.Done()

	return nil
}
