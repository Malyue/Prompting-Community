package storage

import (
	"go.uber.org/zap"
	"prompting/internal/file/config"
	"sync"
)

type CustomStorage interface {
	// MakeBucket 创建存储桶
	MakeBucket(string) error

	// GetObject 获取存储对象
	GetObject(string, string, int64, int64) ([]byte, error)

	// PutObject 上传存储对象
	PutObject(string, string, *[]byte, string) error
}

type ProStorage struct {
	Mux     *sync.RWMutex
	Storage CustomStorage
}

func InitStorage(conf *config.Config) {
	var storageHandler CustomStorage
	if conf.Local.Enabled {
		storageHandler = NewLocalStorage(conf.Local.RootPath)
		zap.L().Info("当前使用的对象存储：" + "Local")
	} else {
		//TODO 其他的对象存储 oss cos等
	}
}
