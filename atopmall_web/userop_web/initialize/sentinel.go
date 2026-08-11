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
	_, err = flow.LoadRules([]*flow.Rule{
		{
			Resource:               "address-list",
			TokenCalculateStrategy: flow.Direct,
			ControlBehavior:        flow.Reject,
			Threshold:              10,
			StatIntervalInMs:       6000,
		},
		{
			Resource:               "address-create",
			TokenCalculateStrategy: flow.Direct,
			ControlBehavior:        flow.Reject,
			Threshold:              3,
			StatIntervalInMs:       6000,
		},
		{
			Resource:               "address-delete",
			TokenCalculateStrategy: flow.Direct,
			ControlBehavior:        flow.Reject,
			Threshold:              5,
			StatIntervalInMs:       6000,
		},
		{
			Resource:               "address-update",
			TokenCalculateStrategy: flow.Direct,
			ControlBehavior:        flow.Reject,
			Threshold:              5,
			StatIntervalInMs:       6000,
		},
		{
			Resource:               "fav-list",
			TokenCalculateStrategy: flow.Direct,
			ControlBehavior:        flow.Reject,
			Threshold:              10,
			StatIntervalInMs:       6000,
		},
		{
			Resource:               "fav-create",
			TokenCalculateStrategy: flow.Direct,
			ControlBehavior:        flow.Reject,
			Threshold:              5,
			StatIntervalInMs:       6000,
		},
		{
			Resource:               "fav-delete",
			TokenCalculateStrategy: flow.Direct,
			ControlBehavior:        flow.Reject,
			Threshold:              5,
			StatIntervalInMs:       6000,
		},
		{
			Resource:               "fav-detail",
			TokenCalculateStrategy: flow.Direct,
			ControlBehavior:        flow.Reject,
			Threshold:              10,
			StatIntervalInMs:       6000,
		},
		{
			Resource:               "message-list",
			TokenCalculateStrategy: flow.Direct,
			ControlBehavior:        flow.Reject,
			Threshold:              10,
			StatIntervalInMs:       6000,
		},
		{
			Resource:               "message-create",
			TokenCalculateStrategy: flow.Direct,
			ControlBehavior:        flow.Reject,
			Threshold:              3,
			StatIntervalInMs:       6000,
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
			MinRequestAmount: 5,
			StatIntervalMs:   10000,
			RetryTimeoutMs:   5000,
		},
		{
			Resource:         "address-create",
			Strategy:         circuitbreaker.ErrorRatio,
			Threshold:        0.5,
			MinRequestAmount: 5,
			StatIntervalMs:   10000,
			RetryTimeoutMs:   5000,
		},
		{
			Resource:         "address-delete",
			Strategy:         circuitbreaker.ErrorRatio,
			Threshold:        0.5,
			MinRequestAmount: 5,
			StatIntervalMs:   10000,
			RetryTimeoutMs:   5000,
		},
		{
			Resource:         "address-update",
			Strategy:         circuitbreaker.ErrorRatio,
			Threshold:        0.5,
			MinRequestAmount: 5,
			StatIntervalMs:   10000,
			RetryTimeoutMs:   5000,
		},
		{
			Resource:         "fav-list",
			Strategy:         circuitbreaker.ErrorRatio,
			Threshold:        0.5,
			MinRequestAmount: 5,
			StatIntervalMs:   10000,
			RetryTimeoutMs:   5000,
		},
		{
			Resource:         "fav-create",
			Strategy:         circuitbreaker.ErrorRatio,
			Threshold:        0.5,
			MinRequestAmount: 5,
			StatIntervalMs:   10000,
			RetryTimeoutMs:   5000,
		},
		{
			Resource:         "fav-delete",
			Strategy:         circuitbreaker.ErrorRatio,
			Threshold:        0.5,
			MinRequestAmount: 5,
			StatIntervalMs:   10000,
			RetryTimeoutMs:   5000,
		},
		{
			Resource:         "fav-detail",
			Strategy:         circuitbreaker.ErrorRatio,
			Threshold:        0.5,
			MinRequestAmount: 5,
			StatIntervalMs:   10000,
			RetryTimeoutMs:   5000,
		},
		{
			Resource:         "message-list",
			Strategy:         circuitbreaker.ErrorRatio,
			Threshold:        0.5,
			MinRequestAmount: 5,
			StatIntervalMs:   10000,
			RetryTimeoutMs:   5000,
		},
		{
			Resource:         "message-create",
			Strategy:         circuitbreaker.ErrorRatio,
			Threshold:        0.5,
			MinRequestAmount: 5,
			StatIntervalMs:   10000,
			RetryTimeoutMs:   5000,
		},
	})
	if err != nil {
		log.Fatalf("配置熔断规则异常:%v", err)
	}
}
