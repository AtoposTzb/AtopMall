package initialize

import (
	"fmt"

	"github.com/fsnotify/fsnotify"
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

	if err := v.Unmarshal(global.ServerConfig); err != nil {
		panic(err)
	}
	zap.S().Infof("配置信息 %s", global.ServerConfig.Name)
	//viper的功能，动态监控变化
	v.WatchConfig()
	v.OnConfigChange(func(in fsnotify.Event) {
		zap.S().Infof("配置文件 %s 变化", in.Name)
		_ = v.ReadInConfig()
		_ = v.Unmarshal(global.ServerConfig)
		zap.S().Infof("配置信息 %s", global.ServerConfig.Name)
	})
}
