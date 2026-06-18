package initialize

import "go.uber.org/zap"

func LoggerInit() {
	//全局的logger
	logger, _ := zap.NewDevelopment()
	defer logger.Sync()
	zap.ReplaceGlobals(logger)
}

/*
	在zap实例中，可以通过以下方式来安全的使用logger
		1.S() 可以获取一个全局的suger实例，可以让我们自己设置一个全局的logger
			我们在main函数中初始化一个全局的logger，后续的代码中就可以直接使用zap.S()来获取全局的suger实例了
		2.日志级别:从低到高分别是:debug ,info ,warn ,error ,dpanic ,panic ,fatal
		3.S()和L(),提供了一个安全的全局访问logger的方法
*/
