package service

import (
	"context"
	"danmu-http/internal/model"
	"danmu-http/internal/validate"
	"danmu-http/logger"
	"danmu-http/middleware"
	"danmu-http/utils"
	"errors"
	"time"
)

const (
	DefaultPassword = "123456"
)

type AuthService interface {
	Register(ctx context.Context, req *validate.RegisterRequest) error
	Login(ctx context.Context, req *validate.LoginRequest) (string, error)
	UpdateSelf(ctx context.Context, req *validate.UpdateSelfRequest) error
	ResetPassword(ctx context.Context, userID string) error
	ListAll(ctx context.Context) ([]*model.Auth, error)
	Delete(ctx context.Context, userID string) error
	GetSelf(ctx context.Context) (*model.Auth, error)
}

type authService struct{}

func NewAuthService() AuthService {
	return &authService{}
}

func (s *authService) Register(ctx context.Context, req *validate.RegisterRequest) error {
	actor, err := middleware.GetAuthFromContext(ctx)
	if err != nil {
		logger.Error().Err(err).Msg("get auth from context failed")
		return err
	}
	logger.Info().
		Str("operator", actor.Email).
		Str("new_email", req.Email).
		Msg("registering new user")

	exists, err := model.IsEmailExists(req.Email)
	if err != nil {
		logger.Error().Err(err).Msg("check email exists failed")
		return err
	}
	if exists {
		return errors.New("email already exists")
	}

	if req.Password == "" {
		req.Password = DefaultPassword
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		logger.Error().Err(err).Msg("hash password failed")
		return err
	}

	if req.Role == "" {
		req.Role = "Guest"
	}

	auth := &model.Auth{
		Email:    req.Email,
		Password: hashedPassword,
		Name:     req.Name,
		Role:     req.Role,
	}

	if err := auth.Insert(); err != nil {
		logger.Error().Err(err).Msg("create auth failed")
		return err
	}

	return nil
}

func (s *authService) Login(ctx context.Context, req *validate.LoginRequest) (string, error) {
	auth, err := model.GetAuthByEmail(req.Email)
	if err != nil {
		logger.Error().Err(err).Msg("get auth by email failed")
		return "", err
	}

	if !utils.CheckPassword(req.Password, auth.Password) {
		return "", errors.New("invalid password")
	}

	token, err := utils.GenerateToken(auth.ID, auth.Email, auth.Name, auth.Role)
	if err != nil {
		logger.Error().Err(err).Msg("generate token failed")
		return "", err
	}
	return token, nil
}

func (s *authService) UpdateSelf(ctx context.Context, req *validate.UpdateSelfRequest) error {
	auth, err := middleware.GetAuthFromContext(ctx)
	if err != nil {
		logger.Error().Err(err).Msg("get auth from context failed")
		return err
	}
	logger.Info().
		Str("operator", auth.Email).
		Str("user_id", auth.ID).
		Msg("updating user profile")
	if req.Password != "" {
		hashedPassword, err := utils.HashPassword(req.Password)
		if err != nil {
			logger.Error().Err(err).Msg("hash password failed")
			return err
		}
		req.Password = hashedPassword
	}

	updateAuth := &model.Auth{
		ID:        auth.ID,
		Email:     req.Email,
		Password:  req.Password,
		Name:      req.Name,
		UpdatedBy: auth.ID,
		UpdatedAt: time.Now().Unix(),
	}

	if err := updateAuth.Update(); err != nil {
		logger.Error().Err(err).Msg("update auth failed")
		return err
	}

	return nil
}

func (s *authService) ResetPassword(ctx context.Context, userID string) error {
	auth, err := middleware.GetAuthFromContext(ctx)
	if err != nil {
		logger.Error().Err(err).Msg("get auth from context failed")
		return err
	}
	logger.Info().
		Str("operator", auth.Email).
		Str("target_user_id", userID).
		Msg("resetting user password")
	oldauth, err := model.GetAuthByID(userID)
	if err != nil {
		logger.Error().Err(err).Msg("get auth by id failed")
		return err
	}

	if oldauth.Role == "admin" {
		return errors.New("cannot reset password for admin")
	}

	initialPassword := DefaultPassword
	hashedPassword, err := utils.HashPassword(initialPassword)
	if err != nil {
		logger.Error().Err(err).Msg("hash password failed")
		return err
	}

	updateAuth := &model.Auth{
		ID:        userID,
		Password:  hashedPassword,
		UpdatedBy: auth.Email,
		UpdatedAt: time.Now().Unix(),
	}

	if err := updateAuth.Update(); err != nil {
		logger.Error().Err(err).Str("auth_email", auth.Email).Str("user_id", userID).Msg("reset password failed")
		return err
	}

	return nil
}

func (s *authService) ListAll(ctx context.Context) ([]*model.Auth, error) {
	auths, err := model.GetAllAuths()
	if err != nil {
		logger.Error().Err(err).Msg("get all auths failed")
		return nil, err
	}

	return auths, nil
}

func (s *authService) Delete(ctx context.Context, userID string) error {
	auth, err := middleware.GetAuthFromContext(ctx)
	if err != nil {
		logger.Error().Err(err).Msg("get auth from context failed")
		return err
	}
	logger.Info().
		Str("operator", auth.Email).
		Str("target_user_id", userID).
		Msg("deleting user")

	auth, err = model.GetAuthByID(userID)
	if err != nil {
		logger.Error().Err(err).Msg("get auth by id failed")
		return err
	}

	if auth.Role == "admin" {
		return errors.New("cannot delete admin")
	}

	deleteAuth := &model.Auth{
		ID: userID,
	}
	if err := deleteAuth.Delete(); err != nil {
		logger.Error().Err(err).Str("auth_email", auth.Email).Str("user_id", userID).Msg("delete auth failed")
		return err
	}

	return nil
}

func (s *authService) GetSelf(ctx context.Context) (*model.Auth, error) {
	auth, err := middleware.GetAuthFromContext(ctx)
	if err != nil {
		logger.Error().Err(err).Msg("get auth from context failed")
		return nil, err
	}

	return model.GetAuthByID(auth.ID)
}
