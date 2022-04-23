package api

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	v1 "github.com/perfectogo/template-service-with-mq/api/handlers/v1"
	"github.com/perfectogo/template-service-with-mq/config"
	"github.com/perfectogo/template-service-with-mq/pkg/logger"
	"github.com/perfectogo/template-service-with-mq/pkg/rabbit"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

func New(cfg config.Config, log logger.Logger, db *sqlx.DB, rmq *rabbit.RMQ) (*gin.Engine, error) {
	if cfg.Environment != "development" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()

	router.Use(gin.Logger(), gin.Recovery()) // Later they will be replaced by custom Logger and Recovery

	handlerV1 := v1.NewHandler(cfg, log, db, rmq)

	router.GET("/ping", handlerV1.Ping)

	rV1 := router.Group("/v1")
	{
		endpointsV1(rV1, handlerV1)
	}

	url := ginSwagger.URL("/swagger/doc.json")
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	return router, nil
}
