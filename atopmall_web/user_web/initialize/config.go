package initialize

import (
	"encoding/json"
	"fmt"

	"github.com/fsnotify/fsnotify"
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"atopmall_web/user_web/global"
)

// 将线上和线下的配置文件隔离
// 设置本地环境变量，viper会自动的获取环境变量情况
func GetEnvInfo(env string) bool {
	viper.AutomaticEnv()
	return viper.GetBool(env)
}

// 设置系统环境变量为ATOPMALL_DEBUG 以后记得重启vscode/traeIDE生效
func ConfigInit() {
	debug := GetEnvInfo(global.Env)
	configFilePrefix := "config"
	configFileName := fmt.Sprintf("%s-pro.yaml", configFilePrefix)
	if debug {
		configFileName = fmt.Sprintf("%s-debug.yaml", configFilePrefix)
	}
	v := viper.New()
	v.SetConfigFile(configFileName)
	if err := v.ReadInConfig(); err != nil {
		panic(err)
	}
	// 解析 nacos 配置 也就是将yaml文件中的配置解析到nNacos配置结构体中,反序列化
	if err := v.Unmarshal(global.NacosConfig); err != nil {
		panic(err)
	}
	zap.S().Infof("配置信息 %s", global.NacosConfig)
	//viper的功能，动态监控变化
	v.WatchConfig()
	v.OnConfigChange(func(in fsnotify.Event) {
		zap.S().Infof("配置文件 %s 变化", in.Name)
		_ = v.ReadInConfig()
		_ = v.Unmarshal(global.NacosConfig)
		zap.S().Infof("配置信息 %s", global.NacosConfig)
	})

	//nacos配置中心配置初始化,参考官网文档
	//1. 配置nacos服务器地址
	// fmt.Println(global.NacosConfig.Host)
	serverConfigs := []constant.ServerConfig{
		{
			IpAddr: global.NacosConfig.Host,
			Port:   global.NacosConfig.Port,
		}}

	//2. 创建clientConfig
	clientConfig := constant.ClientConfig{
		NamespaceId:         global.NacosConfig.NamespaceId, // 如果需要支持多namespace，我们可以创建多个client,它们有不同的NamespaceId。当namespace是public时，此处填空字符串。
		TimeoutMs:           global.NacosConfig.TimeoutMs,
		NotLoadCacheAtStart: global.NacosConfig.NotLoadCacheAtStart,
		LogDir:              global.NacosConfig.LogDir,
		CacheDir:            global.NacosConfig.CacheDir,
		LogLevel:            global.NacosConfig.LogLevel,
	}

	//3. 创建动态配置客户端
	configClient, err := clients.CreateConfigClient(map[string]interface{}{
		"serverConfigs": serverConfigs,
		"clientConfig":  clientConfig,
	})
	if err != nil {
		panic(err)
	}
	//4. 获取配置
	content, err := configClient.GetConfig(vo.ConfigParam{
		DataId: global.NacosConfig.Dataid,
		Group:  global.NacosConfig.Group,
	})
	if err != nil {
		zap.S().Fatalf("获取nacos配置失败 %v", err)
	}
	//5. 解析nacos拿到的配置,此时不在直接读取本地yaml文件了
	err = json.Unmarshal([]byte(content), &global.ServerConfig)
	if err != nil {
		panic(err)
	}
}
