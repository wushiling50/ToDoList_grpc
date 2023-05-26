package discovery

import (
	"errors"
	"fmt"
	"main/ToDoList_grpc/api-gateway/middleware/sen"
	"main/ToDoList_grpc/api-gateway/routes"
	"main/ToDoList_grpc/service"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	sentinel "github.com/alibaba/sentinel-golang/api"
	"github.com/alibaba/sentinel-golang/core/base"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func InitEtcd() {
	//服务发现
	endpoint := []string{viper.GetString("etcd.address")}
	DiscoveryService(endpoint)

	sen.CurrentLimit() //限流
	e, b := sentinel.Entry(sen.ResName, sentinel.WithTrafficType(base.Inbound))
	if b != nil {
		//请求被拒绝
		logrus.Panic("请求被拒绝,BlockErr:", b)
	} else {
		// 请求允许通过,监听转载路由
		go startListen()
		{
			osSignals := make(chan os.Signal, 1)
			signal.Notify(osSignals, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
			s := <-osSignals
			logrus.Println("exit! ", s)

		}
		logrus.Println("网关监听 :4000")
		// 业务结束后调用 Exit
		e.Exit()
		sen.BreakerExit()
	}

}

func startListen() {
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	//连接user
	userConn, err := grpc.Dial(viper.GetString("domain.user.grpcAddress"), opts...)
	if err != nil {
		panic(err)
	}
	var userService service.UserServiceClient

	//连接task
	taskConn, err := grpc.Dial(viper.GetString("domain.task.grpcAddress"), opts...)
	if err != nil {
		panic(err)
	}
	var taskService service.TaskServiceClient

	sen.CircuitBreaker() //熔断
	e, b := sentinel.Entry(sen.Client)
	if b != nil {
		sentinel.TraceError(e, errors.New("biz error"))
	} else {
		userService = service.NewUserServiceClient(userConn)
		taskService = service.NewTaskServiceClient(taskConn)
	}
	sen.E = e

	ginRouter := routes.NewRouter(userService, taskService)

	server := &http.Server{
		Addr:           viper.GetString("server.port"),
		Handler:        ginRouter,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	err = server.ListenAndServe()

	if err != nil {
		panic(fmt.Sprintf("绑定失败,可能端口被占用, err: %v", err))
	}

}
