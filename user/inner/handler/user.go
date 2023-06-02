package handler

import (
	"context"

	"main/pkg/e"
	res "main/pkg/res/user"
	"main/service"
	"main/user/inner/repository"
)

type UserService struct {
	service.UnimplementedUserServiceServer
}

func NewUserService() *UserService {
	return &UserService{}
}

// 用户注册
func (*UserService) UserRegister(ctx context.Context, req *service.UserRequest) (resp *service.UserResponse, err error) {
	var user repository.User
	resp = new(service.UserResponse)
	resp.Code = e.Success
	user, err = user.UserCreate(req)
	if err != nil {
		resp.Code = e.Error
		return resp, err
	}
	resp.UserDetail = res.BulidUser(user)
	return resp, nil
}

// 用户登录
func (*UserService) UserLogin(ctx context.Context, req *service.UserRequest) (resp *service.UserResponse, err error) {
	var user repository.User
	resp = new(service.UserResponse)
	resp.Code = e.Success
	err = user.FindUserInfo(req)
	if err != nil {
		resp.Code = e.Error
		return resp, err
	}
	err = user.CheckPassword(req.Password)
	if err != nil {
		resp.Code = e.Error
		return resp, err
	}
	resp.UserDetail = res.BulidUser(user)
	return resp, nil
}
