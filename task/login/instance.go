package login

import (
	"fmt"
	"main/ToDoList_grpc/pkg/util"
)

type Server struct {
	Name    string `json:"name"`
	Address string `json:"address"`
	Version string `json:"version"` //版本
}

func CreatePrefix(server Server) string {
	if server.Version == "" {
		return fmt.Sprintf("/%s/", server.Name)
	}

	return fmt.Sprintf("/%s/%s/", server.Name, server.Version)
}

// 创建注册路径
func CreateRegisterPath(server Server) string {
	path := fmt.Sprintf("%s%s", CreatePrefix(server), server.Address)
	util.StorePrefix("task", path)
	return path
}
