package service

import (
	"context"
	"danmu-http/internal/model"
	"danmu-http/internal/validate"
	"danmu-http/logger"
	"danmu-http/middleware"
)

type CommonMessageService interface {
	GetCommonMessageWithConditionPage(ctx context.Context, req *validate.CommonMessageQuery) ([]*model.CommonMessage, int64, error)
}

type commonMessageService struct {
}

func NewCommonMessageService() CommonMessageService {
	return &commonMessageService{}
}

func (s *commonMessageService) GetCommonMessageWithConditionPage(ctx context.Context, req *validate.CommonMessageQuery) ([]*model.CommonMessage, int64, error) {
	auth, err := middleware.GetAuthFromContext(ctx)
	if err != nil {
		logger.Error().Err(err).Msg("get auth from context failed")
		return nil, 0, err
	}
	logger.Info().
		Str("operator", auth.Email).
		Interface("req", req).
		Msg("fetching common messages with query")

	messages, count, err := model.GetCommonMessageWithConditionPage(req)
	if err != nil {
		logger.Error().Err(err).Str("operator", auth.Email).Msg("failed to fetch common messages")
		return nil, 0, err
	}

	logger.Info().
		Str("operator", auth.Email).
		Msgf("successfully fetched common messages, count: %d", count)
	return messages, count, nil
}
