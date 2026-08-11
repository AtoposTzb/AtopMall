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
			Resource:               "order-list",
			TokenCalculateStrategy: flow.Direct,
			ControlBehavior:        flow.Reject,
			Threshold:              10,
			StatIntervalInMs:       6000,
		},
		{
			Resource:               "order-detail",
			TokenCalculateStrategy: flow.Direct,
			ControlBehavior:        flow.Reject,
			Threshold:              10,
			StatIntervalInMs:       6000,
		},
		{
			Resource:               "order-create",
			TokenCalculateStrategy: flow.Direct,
			ControlBehavior:        flow.Reject,
			Threshold:              3,
			StatIntervalInMs:       6000,
		},
		{
			Resource:               "cart-list",
			TokenCalculateStrategy: flow.Direct,
			ControlBehavior:        flow.Reject,
			Threshold:              10,
			StatIntervalInMs:       6000,
		},
		{
			Resource:               "cart-create",
			TokenCalculateStrategy: flow.Direct,
			ControlBehavior:        flow.Reject,
			Threshold:              5,
			StatIntervalInMs:       6000,
		},
		{
			Resource:               "cart-delete",
			TokenCalculateStrategy: flow.Direct,
			ControlBehavior:        flow.Reject,
			Threshold:              5,
			StatIntervalInMs:       6000,
		},
		{
			Resource:               "cart-update",
			TokenCalculateStrategy: flow.Direct,
			ControlBehavior:        flow.Reject,
			Threshold:              5,
			StatIntervalInMs:       6000,
		},
		{
			Resource:               "pay-notify",
			TokenCalculateStrategy: flow.Direct,
			ControlBehavior:        flow.Reject,
			Threshold:              10,
			StatIntervalInMs:       6000,
		},
	})
	if err != nil {
		log.Fatalf("配置限流规则异常:%v", err)
	}
	//配置熔断规则
	_, err = circuitbreaker.LoadRules([]*circuitbreaker.Rule{
		{
			Resource:         "order-list",
			Strategy:         circuitbreaker.ErrorRatio,
			Threshold:        0.5,
			MinRequestAmount: 5,
			StatIntervalMs:   10000,
			RetryTimeoutMs:   5000,
		},
		{
			Resource:         "order-detail",
			Strategy:         circuitbreaker.ErrorRatio,
			Threshold:        0.5,
			MinRequestAmount: 5,
			StatIntervalMs:   10000,
			RetryTimeoutMs:   5000,
		},
		{
			Resource:         "order-create",
			Strategy:         circuitbreaker.ErrorRatio,
			Threshold:        0.5,
			MinRequestAmount: 5,
			StatIntervalMs:   10000,
			RetryTimeoutMs:   5000,
		},
		{
			Resource:         "cart-list",
			Strategy:         circuitbreaker.ErrorRatio,
			Threshold:        0.5,
			MinRequestAmount: 5,
			StatIntervalMs:   10000,
			RetryTimeoutMs:   5000,
		},
		{
			Resource:         "cart-create",
			Strategy:         circuitbreaker.ErrorRatio,
			Threshold:        0.5,
			MinRequestAmount: 5,
			StatIntervalMs:   10000,
			RetryTimeoutMs:   5000,
		},
		{
			Resource:         "cart-delete",
			Strategy:         circuitbreaker.ErrorRatio,
			Threshold:        0.5,
			MinRequestAmount: 5,
			StatIntervalMs:   10000,
			RetryTimeoutMs:   5000,
		},
		{
			Resource:         "cart-update",
			Strategy:         circuitbreaker.ErrorRatio,
			Threshold:        0.5,
			MinRequestAmount: 5,
			StatIntervalMs:   10000,
			RetryTimeoutMs:   5000,
		},
		{
			Resource:         "pay-notify",
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
