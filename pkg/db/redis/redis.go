package redis

import (
	"context"
	"encoding/json"
	redis "github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
	"reflect"
	"sync"
	"time"
)

/*
	@Author:Malyue
	@Description:redis初始化
	@CreatedAt:2023/7/6
*/

type Config struct {
	Addr     string `yaml:"addr"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Db       int    `yaml:"db"`
}

var Client *redis.Client

func InitDB(config *Config) error {
	if client != nil {
		return nil
	}
	if config == nil {
		return errors.New("No redis config")
	}
	Client = redis.NewClient(&redis.Options{
		Addr:     config.Addr,
		Password: config.Password,
		Username: config.Username,
		DB:       config.Db,
	})
	callInitHook(Client)
	return nil
}

// InitHook 初始化Redis客户端的钩子函数类型
type InitHook func(client *redis.Client)

var hooks []InitHook
var mutex sync.Mutex

// RegisterInit 注册初始化钩子函数
func RegisterInit(initHook InitHook) {
	mutex.Lock()
	defer mutex.Unlock()
	hooks = append(hooks, initHook)
}

// callInitHook 依次调用注册的初始化狗子函数
func callInitHook(client *redis.Client) {
	mutex.Lock()
	defer mutex.Unlock()
	for _, hook := range hooks {
		hook(client)
	}
}

/*
	Redis的流，是一个日志型数据结构，用于按照时间顺序存储和处理多个事件
*/

// SendEvent 发送事件到Redis Stream，接受一个流名称`stream`和一个事件对象`event`作为参数，在内部调用Redis的`XAdd`方法，事件数据添加到指定流中
func SendEvent(stream string, event interface{}) error {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	data, err := json.Marshal(event)
	if err != nil {
		return err
	}
	add := client.XAdd(ctx, &redis.XAddArgs{
		Stream: stream,
		ID:     "*",
		MaxLen: 50000,
		Approx: true,
		Values: []interface{}{"event", data},
	})
	return add.Err()
}

type EventHandler struct {
	EventType reflect.Type
	Handle    func(eventObj interface{})
}
