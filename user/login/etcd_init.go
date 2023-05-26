package login

import (
	"fmt"
	"main/ToDoList_grpc/service"
	"main/ToDoList_grpc/user/inner/handler"
	"net"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

func InitEtcd() {
	//etcd的地址
	etcdAddress := []string{viper.GetString("etcd.address")}
	// 服务注册
	etcdRegister := NewRegister(etcdAddress, logrus.New())
	grpcAddress := viper.GetString("domain.user.grpcAddress")

	userNode := Server{
		Name:    viper.GetString("domain.user.name"),
		Address: grpcAddress,
	}
	server := grpc.NewServer()

	defer server.Stop()
	// 绑定服务(注册用户服务的服务端)
	service.RegisterUserServiceServer(server, handler.NewUserService())
	//监听
	lis, err := net.Listen("tcp", grpcAddress)
	if err != nil {
		panic(err)
	}
	if _, err := etcdRegister.Register(userNode, 10); err != nil {
		panic(fmt.Sprintf("开启服务失败, err: %v", err))
	}
	logrus.Info("服务开始监听地址： ", grpcAddress)
	if err := server.Serve(lis); err != nil {
		panic(err)
	}
}
