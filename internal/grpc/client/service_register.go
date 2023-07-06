package client

import (
	"github.com/pkg/errors"
	"sync"
)

type initFunc func(config *Service) error

var mutex sync.Mutex
var initFuncRegMap = make(map[string]initFunc)

func RegisterInifFunc(name string, f initFunc) {
	mutex.Lock()
	defer mutex.Unlock()
	initFuncRegMap[name] = f
}

// 注册各api的服务端，以便转发请求
func InitGrpcClient(config *Config) error {
	for _, service := range config.ServiceList {
		fn, o := initFuncRegMap[service.Id]
		if !o {
			return errors.New("The service is not be found")
		}
		// 执行初始化
		if err := fn(service); err != nil {
			return err
		}
	}
	return nil
}
