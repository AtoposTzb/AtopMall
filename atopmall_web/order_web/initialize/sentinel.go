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
	// 读接口：订单列表/详情/购物车，用户高频访问
	// 写接口：下单/购物车操作，重要业务适当放宽
	// pay-notify: 支付回调可能集中到达，需要较高阈值
	_, err = flow.LoadRules([]*flow.Rule{
		{
			Resource:               "order-list",
			TokenCalculateStrategy: flow.Direct,
			ControlBehavior:        flow.Reject,
			Threshold:              30,
			StatIntervalInMs:       1000,
		},
		{
			Resource:               "order-detail",
			TokenCalculateStrategy: flow.Direct,
			ControlBehavior:        flow.Reject,
			Threshold:              30,
			StatIntervalInMs:       1000,
		},
		{
			Resource:               "order-create",
			TokenCalculateStrategy: flow.Direct,
			ControlBehavior:        flow.Reject,
			Threshold:              20,
			StatIntervalInMs:       1000,
		},
		{
			Resource:               "cart-list",
			TokenCalculateStrategy: flow.Direct,
			ControlBehavior:        flow.Reject,
			Threshold:              30,
			StatIntervalInMs:       1000,
		},
		{
			Resource:               "cart-create",
			TokenCalculateStrategy: flow.Direct,
			ControlBehavior:        flow.Reject,
			Threshold:              20,
			StatIntervalInMs:       1000,
		},
		{
			Resource:               "cart-delete",
			TokenCalculateStrategy: flow.Direct,
			ControlBehavior:        flow.Reject,
			Threshold:              20,
			StatIntervalInMs:       1000,
		},
		{
			Resource:               "cart-update",
			TokenCalculateStrategy: flow.Direct,
			ControlBehavior:        flow.Reject,
			Threshold:              20,
			StatIntervalInMs:       1000,
		},
		{
			Resource:               "pay-notify",
			TokenCalculateStrategy: flow.Direct,
			ControlBehavior:        flow.Reject,
			Threshold:              50,
			StatIntervalInMs:       1000,
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
			MinRequestAmount: 20,
			StatIntervalMs:   10000,
			RetryTimeoutMs:   3000,
		},
		{
			Resource:         "order-detail",
			Strategy:         circuitbreaker.ErrorRatio,
			Threshold:        0.5,
			MinRequestAmount: 20,
			StatIntervalMs:   10000,
			RetryTimeoutMs:   3000,
		},
		{
			Resource:         "order-create",
			Strategy:         circuitbreaker.ErrorRatio,
			Threshold:        0.5,
			MinRequestAmount: 10,
			StatIntervalMs:   10000,
			RetryTimeoutMs:   3000,
		},
		{
			Resource:         "cart-list",
			Strategy:         circuitbreaker.ErrorRatio,
			Threshold:        0.5,
			MinRequestAmount: 20,
			StatIntervalMs:   10000,
			RetryTimeoutMs:   3000,
		},
		{
			Resource:         "cart-create",
			Strategy:         circuitbreaker.ErrorRatio,
			Threshold:        0.5,
			MinRequestAmount: 10,
			StatIntervalMs:   10000,
			RetryTimeoutMs:   3000,
		},
		{
			Resource:         "cart-delete",
			Strategy:         circuitbreaker.ErrorRatio,
			Threshold:        0.5,
			MinRequestAmount: 10,
			StatIntervalMs:   10000,
			RetryTimeoutMs:   3000,
		},
		{
			Resource:         "cart-update",
			Strategy:         circuitbreaker.ErrorRatio,
			Threshold:        0.5,
			MinRequestAmount: 10,
			StatIntervalMs:   10000,
			RetryTimeoutMs:   3000,
		},
		{
			Resource:         "pay-notify",
			Strategy:         circuitbreaker.ErrorRatio,
			Threshold:        0.5,
			MinRequestAmount: 20,
			StatIntervalMs:   10000,
			RetryTimeoutMs:   3000,
		},
	})
	if err != nil {
		log.Fatalf("配置熔断规则异常:%v", err)
	}
}
