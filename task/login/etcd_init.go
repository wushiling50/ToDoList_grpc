package login

import (
	"fmt"
	"main/service"
	"main/task/inner/handler"
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
	grpcAddress := viper.GetString("domain.task.grpcAddress")

	taskNode := Server{
		Name:    viper.GetString("domain.task.name"),
		Address: grpcAddress,
	}

	// 绑定服务(服务注册节点)
	server := grpc.NewServer()
	defer server.Stop()
	service.RegisterTaskServiceServer(server, handler.NewTaskService())

	//监听
	lis, err := net.Listen("tcp", grpcAddress)
	if err != nil {
		panic(err)
	}

	// 注册 grpc 服务节点到 etcd 中
	if _, err := etcdRegister.Register(taskNode, 10); err != nil {
		panic(fmt.Sprintf("服务端开启失败, err: %v", err))
	}

	logrus.Info("服务端开始监听地址 ", grpcAddress)
	// 启动 grpc 服务
	if err := server.Serve(lis); err != nil {
		panic(err)
	}
}
