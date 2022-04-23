package todo

import (
	"github.com/jmoiron/sqlx"
	"github.com/perfectogo/template-service-with-mq/config"
	"github.com/perfectogo/template-service-with-mq/pkg/logger"
	"github.com/perfectogo/template-service-with-mq/pkg/rabbit"
	"github.com/perfectogo/template-service-with-mq/storage"
)

type todoEvent struct {
	cfg     config.Config
	log     logger.Logger
	storage storage.StorageInterface
	rmq     *rabbit.RMQ
}

func NewTodoEvent(cfg config.Config, log logger.Logger, db *sqlx.DB, rmq *rabbit.RMQ) *todoEvent {
	return &todoEvent{
		cfg:     cfg,
		log:     log,
		storage: storage.NewStorage(db),
		rmq:     rmq,
	}
}

func (e *todoEvent) RegisterConsumers() {
	e.rmq.AddConsumer(
		"v1.todo.create", // consumerName
		"v1.todo",        // exchangeName
		"v1.todo.create", // queueName
		"v1.todo.create", // routingKey
		e.InsertTodoListener,
	)
}
