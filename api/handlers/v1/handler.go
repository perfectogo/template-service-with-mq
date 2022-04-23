package v1

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/perfectogo/template-service-with-mq/config"
	"github.com/perfectogo/template-service-with-mq/models"
	"github.com/perfectogo/template-service-with-mq/pkg/logger"
	"github.com/perfectogo/template-service-with-mq/pkg/rabbit"
	"github.com/perfectogo/template-service-with-mq/storage"
)

// Handler
type Handler struct {
	cfg       config.Config
	log       logger.Logger
	storagePg storage.StorageInterface
	rmq       *rabbit.RMQ
}

// NewHandler ...
func NewHandler(cfg config.Config, log logger.Logger, db *sqlx.DB, rmq *rabbit.RMQ) *Handler {
	return &Handler{
		cfg:       cfg,
		log:       log,
		storagePg: storage.NewStorage(db),
		rmq:       rmq,
	}
}

// Responses
func (h *Handler) handleSuccessResponse(c *gin.Context, code int, message string, data interface{}) {
	c.JSON(code, models.SuccessModel{
		Code:    fmt.Sprint(code),
		Message: message,
		Data:    data,
	})
}

func (h *Handler) handleErrorResponse(c *gin.Context, code int, message string, err interface{}) {
	h.log.Error(message, logger.Int("code", code), logger.Any("error", err))
	c.JSON(code, models.ErrorModel{
		Code:    fmt.Sprint(code),
		Message: message,
		Error:   err,
	})
}

// QueryParams
func (h *Handler) parseOffsetQueryParam(c *gin.Context) (int, error) {
	return strconv.Atoi(c.DefaultQuery("offset", h.cfg.DefaultOffset))
}

func (h *Handler) parseLimitQueryParam(c *gin.Context) (int, error) {
	return strconv.Atoi(c.DefaultQuery("limit", h.cfg.DefaultLimit))
}
