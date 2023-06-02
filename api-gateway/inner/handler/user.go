package handler

import (
	"context"
	"main/pkg/e"
	"main/pkg/res"
	"main/pkg/util"
	"main/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

// 用户注册
func UserRegister(ginCtx *gin.Context) {
	var userReq service.UserRequest
	UserErrorExist(ginCtx.Bind(&userReq))
	// 从gin.Key中取出服务实例
	userService := ginCtx.Keys["user"].(service.UserServiceClient)
	userResp, err := userService.UserRegister(context.Background(), &userReq) //调用下游的注册服务(Invoke()方法)
	UserErrorExist(err)
	r := res.Response{
		Data:   userResp,
		Status: uint(userResp.Code),
		Msg:    e.GetMsg(uint(userResp.Code)),
	}
	ginCtx.JSON(http.StatusOK, r)
}

// 用户登录
func UserLogin(ginCtx *gin.Context) {
	var userReq service.UserRequest
	UserErrorExist(ginCtx.Bind(&userReq))
	// 从gin.Key中取出服务实例
	userService := ginCtx.Keys["user"].(service.UserServiceClient)
	userResp, err := userService.UserLogin(context.Background(), &userReq) //调用下游的登录服务
	UserErrorExist(err)
	token, err := util.GenerateToken(uint(userResp.UserDetail.Uid))
	if err != nil {
		panic(err)
	}
	r := res.Response{
		Data: res.TokenData{
			User:  userResp.UserDetail,
			Token: token,
		},
		Status: uint(userResp.Code),
		Msg:    e.GetMsg(uint(userResp.Code)),
	}
	ginCtx.JSON(http.StatusOK, r)
}
