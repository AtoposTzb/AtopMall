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
	// sentinel.Init(confPath)
	if err != nil {
		log.Fatalf("初始化sentinel异常:%v", err)
	}
	//配置限流规则
	_, err = flow.LoadRules([]*flow.Rule{
		{
			Resource:               "goods-list",
			TokenCalculateStrategy: flow.Direct, //直接模式
			ControlBehavior:        flow.Reject, //拒绝模式
			Threshold:              3,           //阈值3次/秒
			StatIntervalInMs:       6000,        //统计时间窗口6秒
		},
		{
			Resource:               "create-goods",
			TokenCalculateStrategy: flow.Direct,
			ControlBehavior:        flow.Reject,
			Threshold:              3,
			StatIntervalInMs:       6000,
		},
		{
			Resource:               "goods-detail",
			TokenCalculateStrategy: flow.Direct,
			ControlBehavior:        flow.Reject,
			Threshold:              10,
			StatIntervalInMs:       6000,
		},
		{
			Resource:               "delete-goods",
			TokenCalculateStrategy: flow.Direct,
			ControlBehavior:        flow.Reject,
			Threshold:              3,
			StatIntervalInMs:       6000,
		},
		{
			Resource:               "goods-stocks",
			TokenCalculateStrategy: flow.Direct,
			ControlBehavior:        flow.Reject,
			Threshold:              10,
			StatIntervalInMs:       6000,
		},
		{
			Resource:               "update-goods-status",
			TokenCalculateStrategy: flow.Direct,
			ControlBehavior:        flow.Reject,
			Threshold:              3,
			StatIntervalInMs:       6000,
		},
		{
			Resource:               "update-goods",
			TokenCalculateStrategy: flow.Direct,
			ControlBehavior:        flow.Reject,
			Threshold:              3,
			StatIntervalInMs:       6000,
		},
		{
			Resource:               "category-list",
			TokenCalculateStrategy: flow.Direct,
			ControlBehavior:        flow.Reject,
			Threshold:              10,
			StatIntervalInMs:       6000,
		},
		{
			Resource:               "category-detail",
			TokenCalculateStrategy: flow.Direct,
			ControlBehavior:        flow.Reject,
			Threshold:              10,
			StatIntervalInMs:       6000,
		},
		{
			Resource:               "create-category",
			TokenCalculateStrategy: flow.Direct,
			ControlBehavior:        flow.Reject,
			Threshold:              3,
			StatIntervalInMs:       6000,
		},
		{
			Resource:               "delete-category",
			TokenCalculateStrategy: flow.Direct,
			ControlBehavior:        flow.Reject,
			Threshold:              3,
			StatIntervalInMs:       6000,
		},
		{
			Resource:               "update-category",
			TokenCalculateStrategy: flow.Direct,
			ControlBehavior:        flow.Reject,
			Threshold:              3,
			StatIntervalInMs:       6000,
		},
		{
			Resource:               "brand-list",
			TokenCalculateStrategy: flow.Direct,
			ControlBehavior:        flow.Reject,
			Threshold:              10,
			StatIntervalInMs:       6000,
		},
		{
			Resource:               "create-brand",
			TokenCalculateStrategy: flow.Direct,
			ControlBehavior:        flow.Reject,
			Threshold:              3,
			StatIntervalInMs:       6000,
		},
		{
			Resource:               "delete-brand",
			TokenCalculateStrategy: flow.Direct,
			ControlBehavior:        flow.Reject,
			Threshold:              3,
			StatIntervalInMs:       6000,
		},
		{
			Resource:               "update-brand",
			TokenCalculateStrategy: flow.Direct,
			ControlBehavior:        flow.Reject,
			Threshold:              3,
			StatIntervalInMs:       6000,
		},
		{
			Resource:               "banner-list",
			TokenCalculateStrategy: flow.Direct,
			ControlBehavior:        flow.Reject,
			Threshold:              10,
			StatIntervalInMs:       6000,
		},
		{
			Resource:               "create-banner",
			TokenCalculateStrategy: flow.Direct,
			ControlBehavior:        flow.Reject,
			Threshold:              3,
			StatIntervalInMs:       6000,
		},
		{
			Resource:               "update-banner",
			TokenCalculateStrategy: flow.Direct,
			ControlBehavior:        flow.Reject,
			Threshold:              3,
			StatIntervalInMs:       6000,
		},
		{
			Resource:               "delete-banner",
			TokenCalculateStrategy: flow.Direct,
			ControlBehavior:        flow.Reject,
			Threshold:              3,
			StatIntervalInMs:       6000,
		},
		{
			Resource:               "category-brand-list",
			TokenCalculateStrategy: flow.Direct,
			ControlBehavior:        flow.Reject,
			Threshold:              10,
			StatIntervalInMs:       6000,
		},
		{
			Resource:               "category-brand-all",
			TokenCalculateStrategy: flow.Direct,
			ControlBehavior:        flow.Reject,
			Threshold:              10,
			StatIntervalInMs:       6000,
		},
		{
			Resource:               "create-category-brand",
			TokenCalculateStrategy: flow.Direct,
			ControlBehavior:        flow.Reject,
			Threshold:              3,
			StatIntervalInMs:       6000,
		},
		{
			Resource:               "update-category-brand",
			TokenCalculateStrategy: flow.Direct,
			ControlBehavior:        flow.Reject,
			Threshold:              3,
			StatIntervalInMs:       6000,
		},
		{
			Resource:               "delete-category-brand",
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
			Resource:         "goods-list",
			Strategy:         circuitbreaker.ErrorRatio,
			Threshold:        0.5,
			MinRequestAmount: 5,
			StatIntervalMs:   10000,
			RetryTimeoutMs:   5000,
		},
		{
			Resource:         "create-goods",
			Strategy:         circuitbreaker.ErrorRatio,
			Threshold:        0.5,
			MinRequestAmount: 5,
			StatIntervalMs:   10000,
			RetryTimeoutMs:   5000,
		},
		{
			Resource:         "goods-detail",
			Strategy:         circuitbreaker.ErrorRatio,
			Threshold:        0.5,
			MinRequestAmount: 5,
			StatIntervalMs:   10000,
			RetryTimeoutMs:   5000,
		},
		{
			Resource:         "delete-goods",
			Strategy:         circuitbreaker.ErrorRatio,
			Threshold:        0.5,
			MinRequestAmount: 5,
			StatIntervalMs:   10000,
			RetryTimeoutMs:   5000,
		},
		{
			Resource:         "goods-stocks",
			Strategy:         circuitbreaker.ErrorRatio,
			Threshold:        0.5,
			MinRequestAmount: 5,
			StatIntervalMs:   10000,
			RetryTimeoutMs:   5000,
		},
		{
			Resource:         "update-goods-status",
			Strategy:         circuitbreaker.ErrorRatio,
			Threshold:        0.5,
			MinRequestAmount: 5,
			StatIntervalMs:   10000,
			RetryTimeoutMs:   5000,
		},
		{
			Resource:         "update-goods",
			Strategy:         circuitbreaker.ErrorRatio,
			Threshold:        0.5,
			MinRequestAmount: 5,
			StatIntervalMs:   10000,
			RetryTimeoutMs:   5000,
		},
		{
			Resource:         "category-list",
			Strategy:         circuitbreaker.ErrorRatio,
			Threshold:        0.5,
			MinRequestAmount: 5,
			StatIntervalMs:   10000,
			RetryTimeoutMs:   5000,
		},
		{
			Resource:         "category-detail",
			Strategy:         circuitbreaker.ErrorRatio,
			Threshold:        0.5,
			MinRequestAmount: 5,
			StatIntervalMs:   10000,
			RetryTimeoutMs:   5000,
		},
		{
			Resource:         "create-category",
			Strategy:         circuitbreaker.ErrorRatio,
			Threshold:        0.5,
			MinRequestAmount: 5,
			StatIntervalMs:   10000,
			RetryTimeoutMs:   5000,
		},
		{
			Resource:         "delete-category",
			Strategy:         circuitbreaker.ErrorRatio,
			Threshold:        0.5,
			MinRequestAmount: 5,
			StatIntervalMs:   10000,
			RetryTimeoutMs:   5000,
		},
		{
			Resource:         "update-category",
			Strategy:         circuitbreaker.ErrorRatio,
			Threshold:        0.5,
			MinRequestAmount: 5,
			StatIntervalMs:   10000,
			RetryTimeoutMs:   5000,
		},
		{
			Resource:         "brand-list",
			Strategy:         circuitbreaker.ErrorRatio,
			Threshold:        0.5,
			MinRequestAmount: 5,
			StatIntervalMs:   10000,
			RetryTimeoutMs:   5000,
		},
		{
			Resource:         "create-brand",
			Strategy:         circuitbreaker.ErrorRatio,
			Threshold:        0.5,
			MinRequestAmount: 5,
			StatIntervalMs:   10000,
			RetryTimeoutMs:   5000,
		},
		{
			Resource:         "delete-brand",
			Strategy:         circuitbreaker.ErrorRatio,
			Threshold:        0.5,
			MinRequestAmount: 5,
			StatIntervalMs:   10000,
			RetryTimeoutMs:   5000,
		},
		{
			Resource:         "update-brand",
			Strategy:         circuitbreaker.ErrorRatio,
			Threshold:        0.5,
			MinRequestAmount: 5,
			StatIntervalMs:   10000,
			RetryTimeoutMs:   5000,
		},
		{
			Resource:         "banner-list",
			Strategy:         circuitbreaker.ErrorRatio,
			Threshold:        0.5,
			MinRequestAmount: 5,
			StatIntervalMs:   10000,
			RetryTimeoutMs:   5000,
		},
		{
			Resource:         "create-banner",
			Strategy:         circuitbreaker.ErrorRatio,
			Threshold:        0.5,
			MinRequestAmount: 5,
			StatIntervalMs:   10000,
			RetryTimeoutMs:   5000,
		},
		{
			Resource:         "update-banner",
			Strategy:         circuitbreaker.ErrorRatio,
			Threshold:        0.5,
			MinRequestAmount: 5,
			StatIntervalMs:   10000,
			RetryTimeoutMs:   5000,
		},
		{
			Resource:         "delete-banner",
			Strategy:         circuitbreaker.ErrorRatio,
			Threshold:        0.5,
			MinRequestAmount: 5,
			StatIntervalMs:   10000,
			RetryTimeoutMs:   5000,
		},
		{
			Resource:         "category-brand-list",
			Strategy:         circuitbreaker.ErrorRatio,
			Threshold:        0.5,
			MinRequestAmount: 5,
			StatIntervalMs:   10000,
			RetryTimeoutMs:   5000,
		},
		{
			Resource:         "category-brand-all",
			Strategy:         circuitbreaker.ErrorRatio,
			Threshold:        0.5,
			MinRequestAmount: 5,
			StatIntervalMs:   10000,
			RetryTimeoutMs:   5000,
		},
		{
			Resource:         "create-category-brand",
			Strategy:         circuitbreaker.ErrorRatio,
			Threshold:        0.5,
			MinRequestAmount: 5,
			StatIntervalMs:   10000,
			RetryTimeoutMs:   5000,
		},
		{
			Resource:         "update-category-brand",
			Strategy:         circuitbreaker.ErrorRatio,
			Threshold:        0.5,
			MinRequestAmount: 5,
			StatIntervalMs:   10000,
			RetryTimeoutMs:   5000,
		},
		{
			Resource:         "delete-category-brand",
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
