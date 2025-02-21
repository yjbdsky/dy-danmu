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
		return nil, 0, err
	}

	messages, count, err := model.GetCommonMessageWithConditionPage(req)
	if err != nil {
		logger.Error().
			Err(err).
			Str("operator", auth.Email).
			Interface("query", req).
			Msg("failed to fetch common messages")
		return nil, 0, err
	}

	return messages, count, nil
}
