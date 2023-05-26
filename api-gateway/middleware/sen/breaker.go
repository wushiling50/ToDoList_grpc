package sen

import (
	"log"

	sentinel "github.com/alibaba/sentinel-golang/api"
	"github.com/alibaba/sentinel-golang/core/base"
	"github.com/alibaba/sentinel-golang/core/circuitbreaker"
	"github.com/alibaba/sentinel-golang/util"
	"github.com/sirupsen/logrus"
)

const Client = "client"

type stateChangeTestListener struct {
}

func (s *stateChangeTestListener) OnTransformToClosed(prev circuitbreaker.State, rule circuitbreaker.Rule) {
	logrus.Printf("rule.steategy: %+v, From %s to Closed, time: %d\n", rule.Strategy, prev.String(), util.CurrentTimeMillis())
}

func (s *stateChangeTestListener) OnTransformToOpen(prev circuitbreaker.State, rule circuitbreaker.Rule, snapshot interface{}) {
	logrus.Printf("rule.steategy: %+v, From %s to Open, snapshot: %d, time: %d\n", rule.Strategy, prev.String(), snapshot, util.CurrentTimeMillis())
}

func (s *stateChangeTestListener) OnTransformToHalfOpen(prev circuitbreaker.State, rule circuitbreaker.Rule) {
	logrus.Printf("rule.steategy: %+v, From %s to Half-Open, time: %d\n", rule.Strategy, prev.String(), util.CurrentTimeMillis())
}

func CircuitBreaker() {
	err := sentinel.InitDefault() //初始化
	if err != nil {
		logrus.Fatal("sentinel 初始化失败， err:", err)
	}

	circuitbreaker.RegisterStateChangeListeners(&stateChangeTestListener{}) //注册熔断器监听器

	_, err = circuitbreaker.LoadRules([]*circuitbreaker.Rule{
		// Statistic time span=5s, recoveryTimeout=3s, maxErrorCount=50
		{
			Resource:         "client",
			Strategy:         circuitbreaker.SlowRequestRatio, // 设置熔断断策略为慢调用
			RetryTimeoutMs:   3000,                            // 熔断触发后持续的时间（单位为 ms）(3s)
			MinRequestAmount: 3,                               // 静默数量，对资源的访问小于静默数，熔断器处于静默状态
			StatIntervalMs:   3000,                            // 熔断器的统计周期（单位为 ms）(5s)

			MaxAllowedRtMs: 500, // 慢调用响应时间阈值，RT大于该值的请求判断为慢响应（单位为 ms）
			Threshold:      0.3, // 错误数的阈值
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	logrus.Info("[CircuitBreaker ErrorCount] Sentinel Go circuit breaking  is running.")
}

var E *base.SentinelEntry

func BreakerExit() {
	E.Exit()
}
