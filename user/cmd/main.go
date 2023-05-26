package main

import (
	"main/ToDoList_grpc/conf"
	"main/ToDoList_grpc/user/inner/repository"
	"main/ToDoList_grpc/user/login"
)

func main() {
	conf.InitConfig()
	repository.InitDB()
	login.InitEtcd()
}
