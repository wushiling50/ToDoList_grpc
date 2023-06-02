package main

import (
	"main/conf"
	"main/user/inner/repository"
	"main/user/login"
)

func main() {
	conf.InitConfig()
	repository.InitDB()
	login.InitEtcd()
}
