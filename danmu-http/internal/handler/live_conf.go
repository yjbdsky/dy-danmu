package handler

import (
	"danmu-http/internal/app"
	"danmu-http/internal/service"
	"danmu-http/internal/validate"
	"danmu-http/logger"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type LiveConfHandler struct {
	service service.LiveConfService
}

func NewLiveConfHandler(s service.LiveConfService) *LiveConfHandler {
	return &LiveConfHandler{service: s}
}

func (h *LiveConfHandler) Create(c *gin.Context) {
	var req validate.LiveConfAddRequest
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

	if err := h.service.AddLiveConf(c.Request.Context(), &req); err != nil {
		logger.Error().Err(err).Interface("request", req).Msg("create live conf failed")
		app.NewGin(c).Response(http.StatusInternalServerError, app.ERROR, nil)
		return
	}

	app.NewGin(c).Response(http.StatusOK, app.SUCCESS, nil)
}

func (h *LiveConfHandler) Update(c *gin.Context) {
	var req validate.LiveConfUpdateRequest
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

	if err := h.service.UpdateLiveConf(c.Request.Context(), &req); err != nil {
		logger.Error().Err(err).Interface("request", req).Msg("update live conf failed")
		app.NewGin(c).Response(http.StatusInternalServerError, app.ERROR, nil)
		return
	}

	app.NewGin(c).Response(http.StatusOK, app.SUCCESS, nil)
}

func (h *LiveConfHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		logger.Error().Err(err).Str("id", idStr).Msg("invalid id")
		app.NewGin(c).Response(http.StatusBadRequest, app.InvalidParams, nil)
		return
	}

	if err := h.service.DeleteLiveConf(c.Request.Context(), id); err != nil {
		logger.Error().Err(err).Int64("id", id).Msg("delete live conf failed")
		app.NewGin(c).Response(http.StatusInternalServerError, app.ERROR, nil)
		return
	}

	app.NewGin(c).Response(http.StatusOK, app.SUCCESS, nil)
}

func (h *LiveConfHandler) Get(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		logger.Error().Err(err).Str("id", idStr).Msg("invalid id")
		app.NewGin(c).Response(http.StatusBadRequest, app.InvalidParams, nil)
		return
	}

	conf, err := h.service.GetLiveConfById(c.Request.Context(), id)
	if err != nil {
		logger.Error().Err(err).Int64("id", id).Msg("get live conf failed")
		app.NewGin(c).Response(http.StatusInternalServerError, app.ERROR, nil)
		return
	}

	app.NewGin(c).Response(http.StatusOK, app.SUCCESS, conf)
}

func (h *LiveConfHandler) List(c *gin.Context) {

	confs, err := h.service.ListLiveConf(c.Request.Context())
	if err != nil {
		logger.Error().Err(err).Msg("list live conf failed")
		app.NewGin(c).Response(http.StatusInternalServerError, app.ERROR, nil)
		return
	}

	app.NewGin(c).Response(http.StatusOK, app.SUCCESS, gin.H{
		"list": confs,
	})
}
