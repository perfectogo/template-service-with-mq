package todo

import (
	"encoding/json"
	"fmt"

	"github.com/perfectogo/template-service-with-mq/models"
	"github.com/perfectogo/template-service-with-mq/pkg/logger"
	"github.com/streadway/amqp"
)

func (e *todoEvent) InsertTodoListener(delivery amqp.Delivery) (resp models.SocketResponse) {
	var entity models.CrUpTodo

	//
	if err := json.Unmarshal(delivery.Body, &entity); err != nil {
		e.log.Error("unmarshal error", logger.Any("[]byte", delivery.Body), logger.Any("error", err))
		return models.SocketResponse{}
	}

	//
	todo, err := e.storage.Todo().InsertTodo(entity)
	if err != nil {
		e.log.Error("storage error", logger.Any("entity", entity), logger.Any("error", err))

		return resp
	}

	e.log.Info("todo has been created", logger.Any("entity", entity), logger.Any("res", todo))

	//
	timetable := amqp.Publishing{
		ContentType:   "application/json",
		DeliveryMode:  amqp.Persistent,
		CorrelationId: delivery.CorrelationId,
		Body:          delivery.Body,
	}

	//
	err = e.rmq.Push("v1.todo", "v1.todo.timetable.create-update", timetable)
	if err != nil {
		e.log.Error("publish error", logger.Any("error", err))
		return resp
	}
	fmt.Println(resp)
	return resp
}

//
func (e *todoEvent) UpdateTodoListener(delivery amqp.Delivery) (resp models.SocketResponse) {
	var entity models.CrUpTodo

	//
	if err := json.Unmarshal(delivery.Body, &entity); err != nil {
		e.log.Error("unmarshal error", logger.Any("[]byte", delivery.Body), logger.Any("error", err))
		return models.SocketResponse{}
	}

	//
	todo, err := e.storage.Todo().UpdateTodo(entity)
	if err != nil {
		e.log.Error("storage error", logger.Any("entity", entity), logger.Any("error", err))

		return resp
	}

	e.log.Info("todo has been updated", logger.Any("entity", entity), logger.Any("res", todo))

	//
	timetable := amqp.Publishing{
		ContentType:   "application/json",
		DeliveryMode:  amqp.Persistent,
		CorrelationId: delivery.CorrelationId,
		Body:          delivery.Body,
	}

	//
	err = e.rmq.Push("v1.todo", "v1.todo.timetable.create-update", timetable)
	if err != nil {
		e.log.Error("publish error", logger.Any("error", err))
		return resp
	}
	fmt.Println(resp)
	return resp
}

func (e *todoEvent) DeleteTodoListener(delivery amqp.Delivery) (resp models.SocketResponse) {
	var entity models.ReqById

	//
	if err := json.Unmarshal(delivery.Body, &entity); err != nil {
		e.log.Error("unmarshal error", logger.Any("[]byte", delivery.Body), logger.Any("error", err))
		return models.SocketResponse{}
	}

	//
	err := e.storage.Todo().DeleteTodo(entity)
	if err != nil {
		e.log.Error("storage error", logger.Any("entity", entity), logger.Any("error", err))

		return resp
	}

	e.log.Info("todo has been deleted", logger.Any("entity", entity))

	//
	timetable := amqp.Publishing{
		ContentType:   "application/json",
		DeliveryMode:  amqp.Persistent,
		CorrelationId: delivery.CorrelationId,
		Body:          delivery.Body,
	}

	//
	err = e.rmq.Push("v1.todo", "v1.todo.timetable.create-update", timetable)
	if err != nil {
		e.log.Error("publish error", logger.Any("error", err))
		return resp
	}
	fmt.Println(resp)
	return resp
}
