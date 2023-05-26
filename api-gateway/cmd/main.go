package main

import (
	"main/ToDoList_grpc/api-gateway/discovery"
	"main/ToDoList_grpc/conf"
)

func main() {
	conf.InitConfig()
	discovery.InitEtcd()
}
