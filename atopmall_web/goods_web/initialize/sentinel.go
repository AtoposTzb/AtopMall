package initialize

import (
	"log"

	sentinel "github.com/alibaba/sentinel-golang/api"
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
			//6秒内请求次数超过3次，拒绝请求
		},
	})
	if err != nil {
		log.Fatalf("配置限流规则异常:%v", err)
	}
}
