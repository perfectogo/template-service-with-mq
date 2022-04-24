package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/perfectogo/template-service-with-mq/models"
)

func (h *Handler) Ping(ctx *gin.Context) {
	h.handleSuccessResponse(ctx, 200, "ok", "Pong")
}

//
func (h *Handler) GetTodo(ctx *gin.Context) {
	id := ctx.Param("id")
	todo, err := h.storagePg.Todo().SelectTodo(id)
	if err != nil {
		h.handleErrorResponse(ctx, 400, "bad request", err.Error())
		return
	}
	h.handleSuccessResponse(ctx, 200, "ok", todo)
}

//
func (h *Handler) GetTodoList(ctx *gin.Context) {
	var (
		queryParam models.TodoQueryParamModel
		err        error
	)

	queryParam.Search = ctx.DefaultQuery("search", "")
	queryParam.Order = ctx.DefaultQuery("order", "")
	queryParam.BranchID = ctx.DefaultQuery("branch_id", "")
	queryParam.Arrangement = ctx.DefaultQuery("arrangement", "")
	if len(ctx.Query("offset")) > 0 {
		queryParam.Offset, err = h.parseOffsetQueryParam(ctx)
		if err != nil {
			h.handleErrorResponse(ctx, 400, "wrong offset input", err)
			return
		}
	}

	if len(ctx.Query("limit")) > 0 {
		queryParam.Limit, err = h.parseLimitQueryParam(ctx)
		if err != nil {
			h.handleErrorResponse(ctx, 400, "wrong limit input", err)
			return
		}
	}

	todo, err := h.storagePg.Todo().SelectAllTodo(queryParam)

	if err != nil {
		h.handleErrorResponse(ctx, 400, "bad request", err)
		return
	}

	h.handleSuccessResponse(ctx, 200, "ok", todo)
}
