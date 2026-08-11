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
			Resource:               "user-list",
			TokenCalculateStrategy: flow.Direct,
			ControlBehavior:        flow.Reject,
			Threshold:              5,
			StatIntervalInMs:       6000,
		},
		{
			Resource:               "password-login",
			TokenCalculateStrategy: flow.Direct,
			ControlBehavior:        flow.Reject,
			Threshold:              3,
			StatIntervalInMs:       6000,
		},
		{
			Resource:               "register",
			TokenCalculateStrategy: flow.Direct,
			ControlBehavior:        flow.Reject,
			Threshold:              3,
			StatIntervalInMs:       6000,
		},
		{
			Resource:               "user-detail",
			TokenCalculateStrategy: flow.Direct,
			ControlBehavior:        flow.Reject,
			Threshold:              10,
			StatIntervalInMs:       6000,
		},
		{
			Resource:               "update-user",
			TokenCalculateStrategy: flow.Direct,
			ControlBehavior:        flow.Reject,
			Threshold:              3,
			StatIntervalInMs:       6000,
		},
		{
			Resource:               "send-code",
			TokenCalculateStrategy: flow.Direct,
			ControlBehavior:        flow.Reject,
			Threshold:              1,
			StatIntervalInMs:       60000,
		},
		{
			Resource:               "captcha",
			TokenCalculateStrategy: flow.Direct,
			ControlBehavior:        flow.Reject,
			Threshold:              5,
			StatIntervalInMs:       6000,
		},
	})
	if err != nil {
		log.Fatalf("配置限流规则异常:%v", err)
	}
	//配置熔断规则
	_, err = circuitbreaker.LoadRules([]*circuitbreaker.Rule{
		{
			Resource:         "user-list",
			Strategy:         circuitbreaker.ErrorRatio,
			Threshold:        0.5,
			MinRequestAmount: 5,
			StatIntervalMs:   10000,
			RetryTimeoutMs:   5000,
		},
		{
			Resource:         "password-login",
			Strategy:         circuitbreaker.ErrorRatio,
			Threshold:        0.5,
			MinRequestAmount: 5,
			StatIntervalMs:   10000,
			RetryTimeoutMs:   5000,
		},
		{
			Resource:         "register",
			Strategy:         circuitbreaker.ErrorRatio,
			Threshold:        0.5,
			MinRequestAmount: 5,
			StatIntervalMs:   10000,
			RetryTimeoutMs:   5000,
		},
		{
			Resource:         "user-detail",
			Strategy:         circuitbreaker.ErrorRatio,
			Threshold:        0.5,
			MinRequestAmount: 5,
			StatIntervalMs:   10000,
			RetryTimeoutMs:   5000,
		},
		{
			Resource:         "update-user",
			Strategy:         circuitbreaker.ErrorRatio,
			Threshold:        0.5,
			MinRequestAmount: 5,
			StatIntervalMs:   10000,
			RetryTimeoutMs:   5000,
		},
		{
			Resource:         "send-code",
			Strategy:         circuitbreaker.ErrorRatio,
			Threshold:        0.5,
			MinRequestAmount: 5,
			StatIntervalMs:   10000,
			RetryTimeoutMs:   5000,
		},
		{
			Resource:         "captcha",
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
