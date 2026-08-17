package initialize

import (
	"log"

	sentinel "github.com/alibaba/sentinel-golang/api"
	"github.com/alibaba/sentinel-golang/core/circuitbreaker"
	"github.com/alibaba/sentinel-golang/core/flow"
)

func SentinelInit() {
	// 初始化Sentinel 限流
	err := sentinel.InitDefault()
	if err != nil {
		log.Fatalf("初始化sentinel异常:%v", err)
	}
	//配置限流规则
	// oss-token: 用户上传图片前获取token，前端操作频繁
	// oss-callback: OSS回调，正常流量
	// oss-cleanup: 定时清理任务，低频即可
	_, err = flow.LoadRules([]*flow.Rule{
		{
			Resource:               "oss-token",
			TokenCalculateStrategy: flow.Direct,
			ControlBehavior:        flow.Reject,
			Threshold:              30,
			StatIntervalInMs:       1000,
		},
		{
			Resource:               "oss-upload",
			TokenCalculateStrategy: flow.Direct,
			ControlBehavior:        flow.Reject,
			Threshold:              10,
			StatIntervalInMs:       1000,
		},
		{
			Resource:               "oss-callback",
			TokenCalculateStrategy: flow.Direct,
			ControlBehavior:        flow.Reject,
			Threshold:              30,
			StatIntervalInMs:       1000,
		},
		{
			Resource:               "oss-cleanup",
			TokenCalculateStrategy: flow.Direct,
			ControlBehavior:        flow.Reject,
			Threshold:              1,
			StatIntervalInMs:       60000,
		},
	})
	if err != nil {
		log.Fatalf("配置限流规则异常:%v", err)
	}
	//配置熔断规则
	_, err = circuitbreaker.LoadRules([]*circuitbreaker.Rule{
		{
			Resource:         "oss-token",
			Strategy:         circuitbreaker.ErrorRatio,
			Threshold:        0.5,
			MinRequestAmount: 20,
			StatIntervalMs:   10000,
			RetryTimeoutMs:   3000,
		},
		{
			Resource:         "oss-callback",
			Strategy:         circuitbreaker.ErrorRatio,
			Threshold:        0.5,
			MinRequestAmount: 20,
			StatIntervalMs:   10000,
			RetryTimeoutMs:   3000,
		},
		{
			Resource:         "oss-cleanup",
			Strategy:         circuitbreaker.ErrorRatio,
			Threshold:        0.5,
			MinRequestAmount: 10,
			StatIntervalMs:   10000,
			RetryTimeoutMs:   3000,
		},
	})
	if err != nil {
		log.Fatalf("配置熔断规则异常:%v", err)
	}
}
