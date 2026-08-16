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
	// 读接口：前端用户高频访问，阈值设置较高保证体验流畅
	// 写接口：管理后台操作，适当限流防止滥用
	_, err = flow.LoadRules([]*flow.Rule{
		{
			Resource:               "goods-list",
			TokenCalculateStrategy: flow.Direct,
			ControlBehavior:        flow.Reject,
			Threshold:              100,  // 100 QPS，商品列表是最高频的读接口
			StatIntervalInMs:       1000, // 1秒统计窗口
		},
		{
			Resource:               "goods-detail",
			TokenCalculateStrategy: flow.Direct,
			ControlBehavior:        flow.Reject,
			Threshold:              50, // 50 QPS
			StatIntervalInMs:       1000,
		},
		{
			Resource:               "goods-stocks",
			TokenCalculateStrategy: flow.Direct,
			ControlBehavior:        flow.Reject,
			Threshold:              30, // 30 QPS
			StatIntervalInMs:       1000,
		},
		{
			Resource:               "category-list",
			TokenCalculateStrategy: flow.Direct,
			ControlBehavior:        flow.Reject,
			Threshold:              50, // 50 QPS，首页加载常调用
			StatIntervalInMs:       1000,
		},
		{
			Resource:               "category-detail",
			TokenCalculateStrategy: flow.Direct,
			ControlBehavior:        flow.Reject,
			Threshold:              30, // 30 QPS
			StatIntervalInMs:       1000,
		},
		{
			Resource:               "brand-list",
			TokenCalculateStrategy: flow.Direct,
			ControlBehavior:        flow.Reject,
			Threshold:              50, // 50 QPS，商品筛选常用
			StatIntervalInMs:       1000,
		},
		{
			Resource:               "banner-list",
			TokenCalculateStrategy: flow.Direct,
			ControlBehavior:        flow.Reject,
			Threshold:              30, // 30 QPS
			StatIntervalInMs:       1000,
		},
		{
			Resource:               "category-brand-list",
			TokenCalculateStrategy: flow.Direct,
			ControlBehavior:        flow.Reject,
			Threshold:              30, // 30 QPS
			StatIntervalInMs:       1000,
		},
		{
			Resource:               "category-brand-all",
			TokenCalculateStrategy: flow.Direct,
			ControlBehavior:        flow.Reject,
			Threshold:              30, // 30 QPS
			StatIntervalInMs:       1000,
		},
		{
			Resource:               "create-goods",
			TokenCalculateStrategy: flow.Direct,
			ControlBehavior:        flow.Reject,
			Threshold:              10, // 10 QPS，管理后台操作
			StatIntervalInMs:       1000,
		},
		{
			Resource:               "delete-goods",
			TokenCalculateStrategy: flow.Direct,
			ControlBehavior:        flow.Reject,
			Threshold:              10,
			StatIntervalInMs:       1000,
		},
		{
			Resource:               "update-goods",
			TokenCalculateStrategy: flow.Direct,
			ControlBehavior:        flow.Reject,
			Threshold:              10,
			StatIntervalInMs:       1000,
		},
		{
			Resource:               "update-goods-status",
			TokenCalculateStrategy: flow.Direct,
			ControlBehavior:        flow.Reject,
			Threshold:              10,
			StatIntervalInMs:       1000,
		},
		{
			Resource:               "create-category",
			TokenCalculateStrategy: flow.Direct,
			ControlBehavior:        flow.Reject,
			Threshold:              10,
			StatIntervalInMs:       1000,
		},
		{
			Resource:               "delete-category",
			TokenCalculateStrategy: flow.Direct,
			ControlBehavior:        flow.Reject,
			Threshold:              10,
			StatIntervalInMs:       1000,
		},
		{
			Resource:               "update-category",
			TokenCalculateStrategy: flow.Direct,
			ControlBehavior:        flow.Reject,
			Threshold:              10,
			StatIntervalInMs:       1000,
		},
		{
			Resource:               "create-brand",
			TokenCalculateStrategy: flow.Direct,
			ControlBehavior:        flow.Reject,
			Threshold:              10,
			StatIntervalInMs:       1000,
		},
		{
			Resource:               "delete-brand",
			TokenCalculateStrategy: flow.Direct,
			ControlBehavior:        flow.Reject,
			Threshold:              10,
			StatIntervalInMs:       1000,
		},
		{
			Resource:               "update-brand",
			TokenCalculateStrategy: flow.Direct,
			ControlBehavior:        flow.Reject,
			Threshold:              10,
			StatIntervalInMs:       1000,
		},
		{
			Resource:               "create-banner",
			TokenCalculateStrategy: flow.Direct,
			ControlBehavior:        flow.Reject,
			Threshold:              10,
			StatIntervalInMs:       1000,
		},
		{
			Resource:               "update-banner",
			TokenCalculateStrategy: flow.Direct,
			ControlBehavior:        flow.Reject,
			Threshold:              10,
			StatIntervalInMs:       1000,
		},
		{
			Resource:               "delete-banner",
			TokenCalculateStrategy: flow.Direct,
			ControlBehavior:        flow.Reject,
			Threshold:              10,
			StatIntervalInMs:       1000,
		},
		{
			Resource:               "create-category-brand",
			TokenCalculateStrategy: flow.Direct,
			ControlBehavior:        flow.Reject,
			Threshold:              10,
			StatIntervalInMs:       1000,
		},
		{
			Resource:               "update-category-brand",
			TokenCalculateStrategy: flow.Direct,
			ControlBehavior:        flow.Reject,
			Threshold:              10,
			StatIntervalInMs:       1000,
		},
		{
			Resource:               "delete-category-brand",
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
	// 读接口：MinRequestAmount 设高避免偶发错误误触发熔断，RetryTimeoutMs 缩短加快恢复
	// 写接口：流量较小，MinRequestAmount 适当降低但仍需一定样本量
	_, err = circuitbreaker.LoadRules([]*circuitbreaker.Rule{
		{
			Resource:         "goods-list",
			Strategy:         circuitbreaker.ErrorRatio,
			Threshold:        0.5,
			MinRequestAmount: 30, // 高频接口，需要足够样本量避免误熔断
			StatIntervalMs:   10000,
			RetryTimeoutMs:   3000, // 3秒后尝试恢复
		},
		{
			Resource:         "goods-detail",
			Strategy:         circuitbreaker.ErrorRatio,
			Threshold:        0.5,
			MinRequestAmount: 20,
			StatIntervalMs:   10000,
			RetryTimeoutMs:   3000,
		},
		{
			Resource:         "goods-stocks",
			Strategy:         circuitbreaker.ErrorRatio,
			Threshold:        0.5,
			MinRequestAmount: 20,
			StatIntervalMs:   10000,
			RetryTimeoutMs:   3000,
		},
		{
			Resource:         "category-list",
			Strategy:         circuitbreaker.ErrorRatio,
			Threshold:        0.5,
			MinRequestAmount: 20,
			StatIntervalMs:   10000,
			RetryTimeoutMs:   3000,
		},
		{
			Resource:         "category-detail",
			Strategy:         circuitbreaker.ErrorRatio,
			Threshold:        0.5,
			MinRequestAmount: 20,
			StatIntervalMs:   10000,
			RetryTimeoutMs:   3000,
		},
		{
			Resource:         "brand-list",
			Strategy:         circuitbreaker.ErrorRatio,
			Threshold:        0.5,
			MinRequestAmount: 20,
			StatIntervalMs:   10000,
			RetryTimeoutMs:   3000,
		},
		{
			Resource:         "banner-list",
			Strategy:         circuitbreaker.ErrorRatio,
			Threshold:        0.5,
			MinRequestAmount: 20,
			StatIntervalMs:   10000,
			RetryTimeoutMs:   3000,
		},
		{
			Resource:         "category-brand-list",
			Strategy:         circuitbreaker.ErrorRatio,
			Threshold:        0.5,
			MinRequestAmount: 20,
			StatIntervalMs:   10000,
			RetryTimeoutMs:   3000,
		},
		{
			Resource:         "category-brand-all",
			Strategy:         circuitbreaker.ErrorRatio,
			Threshold:        0.5,
			MinRequestAmount: 20,
			StatIntervalMs:   10000,
			RetryTimeoutMs:   3000,
		},
		{
			Resource:         "create-goods",
			Strategy:         circuitbreaker.ErrorRatio,
			Threshold:        0.5,
			MinRequestAmount: 10, // 写接口流量较小，降低最小请求数
			StatIntervalMs:   10000,
			RetryTimeoutMs:   3000,
		},
		{
			Resource:         "delete-goods",
			Strategy:         circuitbreaker.ErrorRatio,
			Threshold:        0.5,
			MinRequestAmount: 10,
			StatIntervalMs:   10000,
			RetryTimeoutMs:   3000,
		},
		{
			Resource:         "update-goods",
			Strategy:         circuitbreaker.ErrorRatio,
			Threshold:        0.5,
			MinRequestAmount: 10,
			StatIntervalMs:   10000,
			RetryTimeoutMs:   3000,
		},
		{
			Resource:         "update-goods-status",
			Strategy:         circuitbreaker.ErrorRatio,
			Threshold:        0.5,
			MinRequestAmount: 10,
			StatIntervalMs:   10000,
			RetryTimeoutMs:   3000,
		},
		{
			Resource:         "create-category",
			Strategy:         circuitbreaker.ErrorRatio,
			Threshold:        0.5,
			MinRequestAmount: 10,
			StatIntervalMs:   10000,
			RetryTimeoutMs:   3000,
		},
		{
			Resource:         "delete-category",
			Strategy:         circuitbreaker.ErrorRatio,
			Threshold:        0.5,
			MinRequestAmount: 10,
			StatIntervalMs:   10000,
			RetryTimeoutMs:   3000,
		},
		{
			Resource:         "update-category",
			Strategy:         circuitbreaker.ErrorRatio,
			Threshold:        0.5,
			MinRequestAmount: 10,
			StatIntervalMs:   10000,
			RetryTimeoutMs:   3000,
		},
		{
			Resource:         "create-brand",
			Strategy:         circuitbreaker.ErrorRatio,
			Threshold:        0.5,
			MinRequestAmount: 10,
			StatIntervalMs:   10000,
			RetryTimeoutMs:   3000,
		},
		{
			Resource:         "delete-brand",
			Strategy:         circuitbreaker.ErrorRatio,
			Threshold:        0.5,
			MinRequestAmount: 10,
			StatIntervalMs:   10000,
			RetryTimeoutMs:   3000,
		},
		{
			Resource:         "update-brand",
			Strategy:         circuitbreaker.ErrorRatio,
			Threshold:        0.5,
			MinRequestAmount: 10,
			StatIntervalMs:   10000,
			RetryTimeoutMs:   3000,
		},
		{
			Resource:         "create-banner",
			Strategy:         circuitbreaker.ErrorRatio,
			Threshold:        0.5,
			MinRequestAmount: 10,
			StatIntervalMs:   10000,
			RetryTimeoutMs:   3000,
		},
		{
			Resource:         "update-banner",
			Strategy:         circuitbreaker.ErrorRatio,
			Threshold:        0.5,
			MinRequestAmount: 10,
			StatIntervalMs:   10000,
			RetryTimeoutMs:   3000,
		},
		{
			Resource:         "delete-banner",
			Strategy:         circuitbreaker.ErrorRatio,
			Threshold:        0.5,
			MinRequestAmount: 10,
			StatIntervalMs:   10000,
			RetryTimeoutMs:   3000,
		},
		{
			Resource:         "create-category-brand",
			Strategy:         circuitbreaker.ErrorRatio,
			Threshold:        0.5,
			MinRequestAmount: 10,
			StatIntervalMs:   10000,
			RetryTimeoutMs:   3000,
		},
		{
			Resource:         "update-category-brand",
			Strategy:         circuitbreaker.ErrorRatio,
			Threshold:        0.5,
			MinRequestAmount: 10,
			StatIntervalMs:   10000,
			RetryTimeoutMs:   3000,
		},
		{
			Resource:         "delete-category-brand",
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
