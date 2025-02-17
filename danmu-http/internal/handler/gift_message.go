package handler

import (
	"danmu-http/internal/app"
	"danmu-http/internal/service"
	"danmu-http/internal/validate"
	"danmu-http/logger"
	"net/http"

	"github.com/gin-gonic/gin"
)

type GiftMessageHandler struct {
	service service.GiftMessageService
}

func NewGiftMessageHandler(s service.GiftMessageService) *GiftMessageHandler {
	return &GiftMessageHandler{service: s}
}

func (h *GiftMessageHandler) ListGiftRanking(c *gin.Context) {
	var req validate.ListGiftRankingRequest
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

	result, err := h.service.ListGiftRanking(c.Request.Context(), &req)
	if err != nil {
		logger.Error().Err(err).Interface("request", req).Msg("list gift ranking failed")
		app.NewGin(c).Response(http.StatusInternalServerError, app.ERROR, nil)
		return
	}

	app.NewGin(c).Response(http.StatusOK, app.SUCCESS, result)
}

func (h *GiftMessageHandler) ListToUser(c *gin.Context) {
	roomDisplayId := c.Query("room_display_id")
	if roomDisplayId == "" {
		logger.Error().Msg("room_display_id is empty")
		app.NewGin(c).Response(http.StatusBadRequest, app.InvalidParams, nil)
		return
	}

	result, err := h.service.ListToUser(c.Request.Context(), roomDisplayId)
	if err != nil {
		logger.Error().Err(err).Str("room_display_id", roomDisplayId).Msg("list to user failed")
		app.NewGin(c).Response(http.StatusInternalServerError, app.ERROR, nil)
		return
	}

	app.NewGin(c).Response(http.StatusOK, app.SUCCESS, result)
}

func (h *GiftMessageHandler) ListGiftMessagePageWithCondition(c *gin.Context) {
	var req validate.GiftMessageQuery
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

	messages, total, err := h.service.ListGiftMessagePageWithCondition(c.Request.Context(), &req)
	if err != nil {
		logger.Error().Err(err).Interface("request", req).Msg("list gift messages failed")
		app.NewGin(c).Response(http.StatusInternalServerError, app.ERROR, nil)
		return
	}

	app.NewGin(c).Response(http.StatusOK, app.SUCCESS, gin.H{
		"total": total,
		"list":  messages,
	})
}
