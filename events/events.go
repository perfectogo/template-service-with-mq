package events

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/perfectogo/template-service-with-mq/config"
	"github.com/perfectogo/template-service-with-mq/events/todo"
	"github.com/perfectogo/template-service-with-mq/pkg/logger"
	"github.com/perfectogo/template-service-with-mq/pkg/rabbit"
)

type rabbitServer struct {
	cfg config.Config
	log logger.Logger
	db  *sqlx.DB
	RMQ *rabbit.RMQ
}

func NewRabbitServer(cfg config.Config, log logger.Logger, db *sqlx.DB) (*rabbitServer, error) {
	rmq, err := rabbit.NewRMQ(cfg.RabbitURI, log)
	if err != nil {
		return nil, err
	}

	rmq.AddPublisher("v1.todo") // one publisher is enough for application service

	return &rabbitServer{
		cfg: cfg,
		log: log,
		db:  db,
		RMQ: rmq,
	}, nil
}

// Run ...
func (s *rabbitServer) Run(ctx context.Context) {

	todoEvent := todo.NewTodoEvent(s.cfg, s.log, s.db, s.RMQ)
	todoEvent.RegisterConsumers()

	s.RMQ.RunConsumers(ctx) //
}
