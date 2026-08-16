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
	// 读接口（list/detail）：用户高频访问，阈值设置较高
	// 写接口（create/delete/update）：用户操作，适当限流
	_, err = flow.LoadRules([]*flow.Rule{
		{
			Resource:               "address-list",
			TokenCalculateStrategy: flow.Direct,
			ControlBehavior:        flow.Reject,
			Threshold:              30,
			StatIntervalInMs:       1000,
		},
		{
			Resource:               "address-create",
			TokenCalculateStrategy: flow.Direct,
			ControlBehavior:        flow.Reject,
			Threshold:              10,
			StatIntervalInMs:       1000,
		},
		{
			Resource:               "address-delete",
			TokenCalculateStrategy: flow.Direct,
			ControlBehavior:        flow.Reject,
			Threshold:              10,
			StatIntervalInMs:       1000,
		},
		{
			Resource:               "address-update",
			TokenCalculateStrategy: flow.Direct,
			ControlBehavior:        flow.Reject,
			Threshold:              10,
			StatIntervalInMs:       1000,
		},
		{
			Resource:               "fav-list",
			TokenCalculateStrategy: flow.Direct,
			ControlBehavior:        flow.Reject,
			Threshold:              30,
			StatIntervalInMs:       1000,
		},
		{
			Resource:               "fav-create",
			TokenCalculateStrategy: flow.Direct,
			ControlBehavior:        flow.Reject,
			Threshold:              10,
			StatIntervalInMs:       1000,
		},
		{
			Resource:               "fav-delete",
			TokenCalculateStrategy: flow.Direct,
			ControlBehavior:        flow.Reject,
			Threshold:              10,
			StatIntervalInMs:       1000,
		},
		{
			Resource:               "fav-detail",
			TokenCalculateStrategy: flow.Direct,
			ControlBehavior:        flow.Reject,
			Threshold:              30,
			StatIntervalInMs:       1000,
		},
		{
			Resource:               "message-list",
			TokenCalculateStrategy: flow.Direct,
			ControlBehavior:        flow.Reject,
			Threshold:              30,
			StatIntervalInMs:       1000,
		},
		{
			Resource:               "message-create",
			TokenCalculateStrategy: flow.Direct,
			ControlBehavior:        flow.Reject,
			Threshold:              10,
			StatIntervalInMs:       1000,
		},
	})
	if err != nil {
		log.Fatalf("配置限流规则异常:%v", err)
	}
	//配置熔断规则
	_, err = circuitbreaker.LoadRules([]*circuitbreaker.Rule{
		{
			Resource:         "address-list",
			Strategy:         circuitbreaker.ErrorRatio,
			Threshold:        0.5,
			MinRequestAmount: 20,
			StatIntervalMs:   10000,
			RetryTimeoutMs:   3000,
		},
		{
			Resource:         "address-create",
			Strategy:         circuitbreaker.ErrorRatio,
			Threshold:        0.5,
			MinRequestAmount: 10,
			StatIntervalMs:   10000,
			RetryTimeoutMs:   3000,
		},
		{
			Resource:         "address-delete",
			Strategy:         circuitbreaker.ErrorRatio,
			Threshold:        0.5,
			MinRequestAmount: 10,
			StatIntervalMs:   10000,
			RetryTimeoutMs:   3000,
		},
		{
			Resource:         "address-update",
			Strategy:         circuitbreaker.ErrorRatio,
			Threshold:        0.5,
			MinRequestAmount: 10,
			StatIntervalMs:   10000,
			RetryTimeoutMs:   3000,
		},
		{
			Resource:         "fav-list",
			Strategy:         circuitbreaker.ErrorRatio,
			Threshold:        0.5,
			MinRequestAmount: 20,
			StatIntervalMs:   10000,
			RetryTimeoutMs:   3000,
		},
		{
			Resource:         "fav-create",
			Strategy:         circuitbreaker.ErrorRatio,
			Threshold:        0.5,
			MinRequestAmount: 10,
			StatIntervalMs:   10000,
			RetryTimeoutMs:   3000,
		},
		{
			Resource:         "fav-delete",
			Strategy:         circuitbreaker.ErrorRatio,
			Threshold:        0.5,
			MinRequestAmount: 10,
			StatIntervalMs:   10000,
			RetryTimeoutMs:   3000,
		},
		{
			Resource:         "fav-detail",
			Strategy:         circuitbreaker.ErrorRatio,
			Threshold:        0.5,
			MinRequestAmount: 20,
			StatIntervalMs:   10000,
			RetryTimeoutMs:   3000,
		},
		{
			Resource:         "message-list",
			Strategy:         circuitbreaker.ErrorRatio,
			Threshold:        0.5,
			MinRequestAmount: 20,
			StatIntervalMs:   10000,
			RetryTimeoutMs:   3000,
		},
		{
			Resource:         "message-create",
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
