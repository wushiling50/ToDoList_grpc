package main

import (
	"main/conf"
	"main/task/inner/repository"
	"main/task/login"
)

func main() {
	conf.InitConfig()
	repository.InitDB()
	login.InitEtcd()
}
