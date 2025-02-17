package handler

import (
	"danmu-http/internal/app"
	"danmu-http/internal/service"
	"danmu-http/internal/validate"
	"danmu-http/logger"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	service service.AuthService
}

func NewAuthHandler(s service.AuthService) *AuthHandler {
	return &AuthHandler{service: s}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req validate.RegisterRequest
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

	if err := h.service.Register(c.Request.Context(), &req); err != nil {
		logger.Error().Err(err).Interface("request", req).Msg("register user failed")
		app.NewGin(c).Response(http.StatusInternalServerError, app.ERROR, nil)
		return
	}

	app.NewGin(c).Response(http.StatusOK, app.SUCCESS, nil)
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req validate.LoginRequest
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

	token, err := h.service.Login(c.Request.Context(), &req)
	if err != nil {
		logger.Error().Err(err).Str("email", req.Email).Msg("login failed")
		app.NewGin(c).Response(http.StatusUnauthorized, app.Unauthorized, nil)
		return
	}

	app.NewGin(c).Response(http.StatusOK, app.SUCCESS, gin.H{"token": token})
}

func (h *AuthHandler) UpdateSelf(c *gin.Context) {
	var req validate.UpdateSelfRequest
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

	if err := h.service.UpdateSelf(c.Request.Context(), &req); err != nil {
		logger.Error().Err(err).Interface("request", req).Msg("update self failed")
		app.NewGin(c).Response(http.StatusInternalServerError, app.ERROR, nil)
		return
	}

	app.NewGin(c).Response(http.StatusOK, app.SUCCESS, nil)
}

func (h *AuthHandler) ResetPassword(c *gin.Context) {
	userID := c.Param("id")
	if userID == "" {
		logger.Error().Msg("user id is empty")
		app.NewGin(c).Response(http.StatusBadRequest, app.InvalidParams, nil)
		return
	}

	if err := h.service.ResetPassword(c.Request.Context(), userID); err != nil {
		logger.Error().Err(err).Str("user_id", userID).Msg("reset password failed")
		app.NewGin(c).Response(http.StatusInternalServerError, app.ERROR, nil)
		return
	}

	app.NewGin(c).Response(http.StatusOK, app.SUCCESS, nil)
}

func (h *AuthHandler) ListAll(c *gin.Context) {
	users, err := h.service.ListAll(c.Request.Context())
	if err != nil {
		logger.Error().Err(err).Msg("list users failed")
		app.NewGin(c).Response(http.StatusInternalServerError, app.ERROR, nil)
		return
	}

	app.NewGin(c).Response(http.StatusOK, app.SUCCESS, users)
}

func (h *AuthHandler) Delete(c *gin.Context) {
	userID := c.Param("id")
	if userID == "" {
		logger.Error().Msg("user id is empty")
		app.NewGin(c).Response(http.StatusBadRequest, app.InvalidParams, nil)
		return
	}

	if err := h.service.Delete(c.Request.Context(), userID); err != nil {
		logger.Error().Err(err).Str("user_id", userID).Msg("delete user failed")
		app.NewGin(c).Response(http.StatusInternalServerError, app.ERROR, nil)
		return
	}

	app.NewGin(c).Response(http.StatusOK, app.SUCCESS, nil)
}

func (h *AuthHandler) GetSelf(c *gin.Context) {
	user, err := h.service.GetSelf(c.Request.Context())
	if err != nil {
		logger.Error().Err(err).Str("user_id", user.ID).Msg("get self failed")
		app.NewGin(c).Response(http.StatusInternalServerError, app.ERROR, nil)
		return
	}

	app.NewGin(c).Response(http.StatusOK, app.SUCCESS, user)
}
