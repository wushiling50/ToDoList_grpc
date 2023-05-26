package sen

import (
	sentinel "github.com/alibaba/sentinel-golang/api"
	"github.com/alibaba/sentinel-golang/core/flow"
	"github.com/sirupsen/logrus"
)

const ResName string = "todolist-grpc"

func CurrentLimit() {
	err := sentinel.InitDefault() //使用默认的初始化
	if err != nil {
		logrus.Fatal("sentinel 初始化失败， err:", err)
	}

	//配置限流规则
	_, err = flow.LoadRules([]*flow.Rule{
		{
			Resource:               ResName,
			TokenCalculateStrategy: flow.Direct,
			ControlBehavior:        flow.Throttling,
			Threshold:              10,   //请求间隔为100ms
			StatIntervalInMs:       1000, //统计周期1s
			MaxQueueingTimeMs:      500,  // 最长排队等待时间
		},
	})
	if err != nil {
		logrus.Fatal("限流规则配置失败, err:", err)
	}
}
