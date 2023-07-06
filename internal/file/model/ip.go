package model

import (
	"context"
	"prompting/pkg/db/redis"
)

func GetIpExist(ip string) int64 {
	return redis.Client.Exists(context.Background(), ip).Val()
}

func GetIp(ip string) (string, error) {
	return redis.Client.Get(context.Background(), ip).Result()
}

// 自增工作ID
func IncrWorkId(key string) (int64, error) {
	return redis.Client.Incr(context.Background(), key).Result()
}

func SetWorkIdMapToIp(ip string, workId int64) {
	redis.Client.Set(context.Background(), ip, workId, -1)
}
