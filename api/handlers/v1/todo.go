package v1

import "github.com/gin-gonic/gin"

func (h *Handler) Ping(ctx *gin.Context) {
	h.handleSuccessResponse(ctx, 200, "ok", "Pong")
}

func (h *Handler) GetTodo(ctx *gin.Context) {
	id := ctx.Param("id")
	todo, err := h.storagePg.Todo().SelectTodo(id)
	if err != nil {
		h.handleErrorResponse(ctx, 400, "bad request", err.Error())
		return
	}
	h.handleSuccessResponse(ctx, 200, "ok", todo)
}
