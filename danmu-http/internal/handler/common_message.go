package handler

import (
	"danmu-http/internal/app"
	"danmu-http/internal/service"
	"danmu-http/internal/validate"
	"danmu-http/logger"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CommonMessageHandler struct {
	service service.CommonMessageService
}

func NewCommonMessageHandler(s service.CommonMessageService) *CommonMessageHandler {
	return &CommonMessageHandler{service: s}
}

func (h *CommonMessageHandler) ListPageableWithCondition(c *gin.Context) {
	var req validate.CommonMessageQuery
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error().Err(err).Msg("bind request failed")
		app.NewGin(c).Response(http.StatusBadRequest, app.InvalidParams, nil)
		return
	}

	if err := validate.Struct(req); err != nil {
		logger.Error().Err(err).Interface("request", req).Msg("validate request failed")
		app.NewGin(c).Response(http.StatusBadRequest, app.InvalidParams, err.Error())
		return
	}

	messages, total, err := h.service.GetCommonMessageWithConditionPage(c.Request.Context(), &req)
	if err != nil {
		logger.Error().Err(err).Interface("request", req).Msg("list common messages failed")
		app.NewGin(c).Response(http.StatusInternalServerError, app.ERROR, nil)
		return
	}

	app.NewGin(c).Response(http.StatusOK, app.SUCCESS, gin.H{
		"total": total,
		"list":  messages,
	})
}
