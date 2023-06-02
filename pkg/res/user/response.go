package res

import (
	"main/service"
	"main/user/inner/repository"
)

func BulidUser(item repository.User) *service.UserModel {
	userModel := service.UserModel{
		Uid:      uint32(item.Uid),
		UserName: item.UserName,
	}
	return &userModel
}
