package main

import (
	"flag"
	"go.uber.org/zap"
	"prompting/internal/file/plugins"
	"prompting/internal/file/utils"
	"prompting/internal/gateway/config"
	"prompting/pkg/db/redis"
)

/*
	Author:Malyue
	Description:文件模块
	CreatedAt:2023/6/7
*/

func main() {
	var configFilename = flag.String("config", "config/gateway.yaml", "")
	flag.Parse()

	cfg, err := config.LoadFromConfigFile(*configFilename)
	if err != nil {
		zap.L().Error("[File Read Config Error:]", zap.Error(err))
		panic(0)
	}

	// 初始化插件
	plugins.NewPlugins()
	defer plugins.ClosePlugins()

	// 初始化redis
	err = redis.InitDB(cfg.Redis)
	if err != nil {
		zap.L().Error("[Init Redis Error:]", zap.Error(err))
		return
	}

	// 初始化雪花算法
	utils.InitSnowFlake()

	// 初始化存储

	// 初始化路由

}
