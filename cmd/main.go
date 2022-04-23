package main

import (
	"context"
	"time"

	"github.com/perfectogo/template-service-with-mq/api"
	"github.com/perfectogo/template-service-with-mq/config"
	"github.com/perfectogo/template-service-with-mq/events"
	"github.com/perfectogo/template-service-with-mq/pkg/db"
	"github.com/perfectogo/template-service-with-mq/pkg/logger"
	"golang.org/x/sync/errgroup"
)

func main() {
	// load config
	cfg := config.Load()

	//
	log := logger.New(cfg.App, cfg.LogLevel)

	//
	sqlx, err := db.ConnectToDB(cfg)
	if err != nil {
		log.Panic("error connecting to postgres", logger.Error(err))
	}
	sqlx.SetMaxOpenConns(50)
	sqlx.SetMaxIdleConns(25)
	sqlx.SetConnMaxLifetime(5 * time.Minute)

	//
	rabbitServer, err := events.NewRabbitServer(cfg, log, sqlx)
	if err != nil {
		log.Panic("error on the event server", logger.Error(err))
	}

	//
	apiServer, err := api.New(cfg, log, sqlx, rabbitServer.RMQ)
	if err != nil {
		log.Panic("error on the api server", logger.Error(err))
	}

	//
	group, ctx := errgroup.WithContext(context.Background())

	//
	group.Go(func() error {
		rabbitServer.Run(ctx) // it should run forever if there is any consumer
		log.Panic("event server has finished")
		return nil
	})

	group.Go(func() error {
		err := apiServer.Run(cfg.HTTPPort) // this method will block the calling goroutine indefinitely unless an error happens
		if err != nil {
			panic(err)
		}
		log.Panic("api server has finished")
		return nil
	})

	err = group.Wait()
	if err != nil {
		log.Panic("error on the server", logger.Error(err))
	}
}
