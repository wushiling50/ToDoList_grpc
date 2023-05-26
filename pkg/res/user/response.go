package res

import (
	"main/ToDoList_grpc/service"
	"main/ToDoList_grpc/user/inner/repository"
)

func BulidUser(item repository.User) *service.UserModel {
	userModel := service.UserModel{
		Uid:      uint32(item.Uid),
		UserName: item.UserName,
	}
	return &userModel
}
