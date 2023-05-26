package main

import (
	"main/ToDoList_grpc/conf"
	"main/ToDoList_grpc/task/inner/repository"
	"main/ToDoList_grpc/task/login"
)

func main() {
	conf.InitConfig()
	repository.InitDB()
	login.InitEtcd()
}
