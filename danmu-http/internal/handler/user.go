package handler

import (
	"danmu-http/internal/app"
	"danmu-http/internal/service"
	"danmu-http/internal/validate"
	"danmu-http/logger"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	service service.UserService
}

func NewUserHandler(s service.UserService) *UserHandler {
	return &UserHandler{service: s}
}

func (h *UserHandler) ListAllUsers(c *gin.Context) {
	var req *validate.UserPageRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		logger.Error().Err(err).Msg("bind request failed")
		app.NewGin(c).Response(http.StatusBadRequest, app.InvalidParams, nil)
		return
	}
	users, total, err := h.service.ListAllUsers(c.Request.Context(), req)
	if err != nil {
		logger.Error().Err(err).Msg("list all users failed")
		app.NewGin(c).Response(http.StatusInternalServerError, app.ERROR, nil)
		return
	}

	app.NewGin(c).Response(http.StatusOK, app.SUCCESS, gin.H{
		"total": total,
		"list":  users,
	})
}

func (h *UserHandler) SearchUser(c *gin.Context) {
	var req *validate.UserSearchRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		logger.Error().Err(err).Msg("bind request failed")
		app.NewGin(c).Response(http.StatusBadRequest, app.InvalidParams, nil)
		return
	}

	users, total, err := h.service.SearchUser(c.Request.Context(), req)
	if err != nil {
		logger.Error().Err(err).Str("keyword", req.Keyword).Msg("search user failed")
		app.NewGin(c).Response(http.StatusInternalServerError, app.ERROR, nil)
		return
	}

	app.NewGin(c).Response(http.StatusOK, app.SUCCESS, gin.H{
		"total": total,
		"list":  users,
	})
}
