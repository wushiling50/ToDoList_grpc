package main

import (
	"main/api-gateway/discovery"
	"main/conf"
)

func main() {
	conf.InitConfig()
	discovery.InitEtcd()
}
