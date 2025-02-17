package service

import (
	"context"
	"danmu-http/internal/model"
	"danmu-http/internal/validate"
)

type UserService interface {
	ListAllUsers(ctx context.Context, req *validate.UserPageRequest) ([]*model.User, int64, error)
	SearchUser(ctx context.Context, req *validate.UserSearchRequest) ([]*model.User, int64, error)
}

type userService struct {
}

func NewUserService() UserService {
	return &userService{}
}

func (s *userService) ListAllUsers(ctx context.Context, req *validate.UserPageRequest) ([]*model.User, int64, error) {
	return model.GetAllUsersPage(req.Page, req.PageSize)
}

func (s *userService) SearchUser(ctx context.Context, req *validate.UserSearchRequest) ([]*model.User, int64, error) {
	return model.SearchUser(req.Page, req.PageSize, req.Keyword)
}
