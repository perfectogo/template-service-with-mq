package api

import (
	"github.com/gin-gonic/gin"
	v1 "github.com/perfectogo/template-service-with-mq/api/handlers/v1"
)

func endpointsV1(r *gin.RouterGroup, h *v1.Handler) {
	r.GET("/todo/:id", h.GetTodo)
}
