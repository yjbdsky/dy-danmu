package service

import (
	"context"
	"danmu-http/internal/model"
	"danmu-http/internal/validate"
	"danmu-http/logger"
	"danmu-http/middleware"
	"danmu-http/rpc"
	api "danmu-http/rpc/proto"
	"errors"
	"time"

	"gorm.io/gorm"
)

type LiveConfService interface {
	ListLiveConf(ctx context.Context) ([]*model.LiveConf, error)
	AddLiveConf(ctx context.Context, req *validate.LiveConfAddRequest) error
	UpdateLiveConf(ctx context.Context, req *validate.LiveConfUpdateRequest) error
	DeleteLiveConf(ctx context.Context, id int64) error
	GetLiveConfById(ctx context.Context, id int64) (*model.LiveConf, error)
}

type liveConfService struct {
}

func NewLiveConfService() LiveConfService {
	return &liveConfService{}
}

func (s *liveConfService) ListLiveConf(ctx context.Context) ([]*model.LiveConf, error) {
	return model.GetAllLiveConf()
}

func (s *liveConfService) AddLiveConf(ctx context.Context, req *validate.LiveConfAddRequest) error {
	auth, err := middleware.GetAuthFromContext(ctx)
	if err != nil {
		return err
	}

	if auth.Name == "" {
		auth.Name = "system"
	}

	logger.Info().
		Str("operator", auth.Email).
		Str("room_id", req.RoomDisplayID).
		Str("name", req.Name).
		Bool("enable", req.Enable).
		Msg("adding new live configuration")

	now := time.Now().Unix()
	liveConf := &model.LiveConf{
		RoomDisplayID: req.RoomDisplayID,
		URL:           req.URL,
		Name:          req.Name,
		Enable:        req.Enable,
		ModifiedBy:    auth.Email,
		CratedBy:      auth.Email,
		ModifiedOn:    now,
		CreatedOn:     now,
	}

	err = model.DB.Transaction(func(tx *gorm.DB) error {
		if err := liveConf.Insert(tx); err != nil {
			logger.Error().
				Err(err).
				Str("operator", auth.Email).
				Str("room_id", req.RoomDisplayID).
				Msg("failed to add live configuration")
			return err
		}

		rpcClient, err := rpc.GetClient()
		if err != nil {
			logger.Error().Err(err).Str("auth_id", auth.ID).Str("auth_name", auth.Name).Msg("get rpc client failed")
			return err
		}

		ctx, cancel := rpc.WithTimeout(5 * time.Second)
		defer cancel()

		rpcReq := &api.LiveConf{
			Id:            liveConf.ID,
			RoomDisplayId: liveConf.RoomDisplayID,
			Url:           liveConf.URL,
			Name:          liveConf.Name,
			Enable:        liveConf.Enable,
		}
		res, err := rpcClient.AddTask(ctx, rpcReq)
		if err != nil {
			logger.Error().Err(err).Str("auth_id", auth.ID).Str("auth_name", auth.Name).Msg("add live conf to rpc failed")
			return err
		}
		if res.Code != 200 {
			logger.Error().Str("auth_id", auth.ID).Str("auth_name", auth.Name).Msg("add live conf to rpc failed")
			return errors.New(res.Message)
		}
		return nil
	})
	return err
}

func (s *liveConfService) UpdateLiveConf(ctx context.Context, req *validate.LiveConfUpdateRequest) error {
	auth, err := middleware.GetAuthFromContext(ctx)
	if err != nil {
		return err
	}

	logger.Info().
		Str("operator", auth.Email).
		Int64("conf_id", req.ID).
		Str("room_id", req.RoomDisplayID).
		Bool("enable", req.Enable).
		Msg("updating live configuration")

	conf, err := model.GetLiveConfById(req.ID)
	if err != nil {
		logger.Error().Err(err).Str("auth_id", auth.ID).Str("auth_name", auth.Name).Msg("get live conf by id failed")
		return err
	}

	conf.RoomDisplayID = req.RoomDisplayID
	conf.URL = req.URL
	conf.Name = req.Name
	conf.Enable = req.Enable
	conf.ModifiedBy = auth.Email
	conf.ModifiedOn = time.Now().Unix()

	err = model.DB.Transaction(func(tx *gorm.DB) error {
		if err := conf.Update(tx); err != nil {
			logger.Error().Err(err).Str("auth_id", auth.ID).Str("auth_name", auth.Name).Msg("update live conf failed")
			return err
		}

		rpcClient, err := rpc.GetClient()
		if err != nil {
			logger.Error().Err(err).Str("auth_id", auth.ID).Str("auth_name", auth.Name).Msg("get rpc client failed")
			return err
		}

		ctx, cancel := rpc.WithTimeout(5 * time.Second)
		defer cancel()

		rpcReq := &api.LiveConf{
			Id:            conf.ID,
			RoomDisplayId: conf.RoomDisplayID,
			Url:           conf.URL,
			Name:          conf.Name,
			Enable:        conf.Enable,
		}
		res, err := rpcClient.UpdateTask(ctx, rpcReq)
		if err != nil {
			logger.Error().Err(err).Str("auth_id", auth.ID).Str("auth_name", auth.Name).Msg("update live conf to rpc failed")
			return err
		}
		if res.Code != 200 {
			logger.Error().Str("auth_id", auth.ID).Str("auth_name", auth.Name).Msg("update live conf to rpc failed")
			return errors.New(res.Message)
		}
		return nil
	})
	return err
}

func (s *liveConfService) DeleteLiveConf(ctx context.Context, id int64) error {
	auth, err := middleware.GetAuthFromContext(ctx)
	if err != nil {
		return err
	}

	logger.Info().
		Str("operator", auth.Email).
		Int64("conf_id", id).
		Msg("deleting live configuration")

	rpcClient, err := rpc.GetClient()
	if err != nil {
		logger.Error().Err(err).Str("auth_id", auth.ID).Str("auth_name", auth.Name).Msg("get rpc client failed")
		return err
	}

	ctx, cancel := rpc.WithTimeout(5 * time.Second)
	defer cancel()

	if _, err := rpcClient.DeleteTask(ctx, &api.TaskID{Id: id}); err != nil {
		logger.Error().Err(err).Str("auth_id", auth.ID).Str("auth_name", auth.Name).Msg("delete live conf to rpc failed")
		return err
	}

	if err := model.DeleteLiveConfById(id); err != nil {
		logger.Error().Err(err).Str("auth_id", auth.ID).Str("auth_name", auth.Name).Msg("delete live conf failed")
		return err
	}
	return nil
}

func (s *liveConfService) GetLiveConfById(ctx context.Context, id int64) (*model.LiveConf, error) {
	return model.GetLiveConfById(id)
}
