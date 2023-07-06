package main

import (
	"flag"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"prompting/internal/gateway/config"
	"prompting/internal/gateway/server"
	"prompting/internal/grpc/client"
	"prompting/pkg/db/mysql"
	"prompting/pkg/db/redis"
	"prompting/pkg/log"
	"syscall"
)

func main() {
	var configFilename = flag.String("config", "config/gateway.yaml", "")
	flag.Parse()
	// 加载配置文件
	cfg, err := config.LoadFromConfigFile(*configFilename)
	if err != nil {
		zap.L().Error("[GateWay Read Config Error:]", zap.Error(err))
		panic(0)
	}

	// 初始化日志,因为gateway是第一个启动的，所以可以在这里实现
	err = log.InitLog(*cfg.Log, cfg.Mode)
	if err != nil {
		zap.L().Error("[Init log error]:", zap.Error(err))
		return
	}

	// 初始化mysql
	err = mysql.InitDB(cfg.Mysql)

	// 初始化redis
	err = redis.InitDB(cfg.Redis)

	// 创建客户端与其他模块通讯
	err = client.InitGrpcClient(cfg.Client)
	if err != nil {
		zap.L().Error("[Init grpc client error]:", zap.Error(err))
		return
	}

	// 创建服务端与客户端通讯
	svr := server.NewServer(cfg.Server)
	zap.L().Info("The server is start")
	err = svr.Run()
	if err != nil {
		zap.L().Error("Server Startup fail", zap.Error(err))
		os.Exit(1)
	}

	var exitFunc = func() {
		zap.L().Info("The server is exit")
		svr.Close()
	}
	go func() {
		sc := make(chan os.Signal, 1)
		signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
		<-sc
		exitFunc()
		zap.L().Info("The server has stop")
	}()
}
