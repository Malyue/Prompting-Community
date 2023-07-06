package utils

import (
	"errors"
	"go.uber.org/zap"
	"prompting/internal/file/model"
	"strconv"
	"sync"
	"time"
)

type Snowflake struct {
	mu            sync.Mutex
	lastTimestamp int64
	workerId      int64
	datacenterId  int64
	sequence      int64
}

var (
	snowFlake *Snowflake
	once      sync.Once
)

const (
	twepoch            = int64(1417937700000) // Unix纪元时间戳
	workerIdBits       = uint(5)              // 机器ID所占位数
	datacenterBits     = uint(5)              // 数据中心ID所占位数
	maxWorkerId        = int64(-1) ^ (int64(-1) << workerIdBits)
	maxDatacenterId    = int64(-1) ^ (int64(-1) << datacenterBits)
	sequenceBits       = uint(12) // 序列号所占位数
	workerIdShift      = sequenceBits
	datacenterIdShift  = sequenceBits + workerIdBits
	timestampLeftShift = sequenceBits + workerIdBits + datacenterBits
	sequenceMask       = int64(-1) ^ (int64(-1) << sequenceBits)
)

// InitSnowFlake 分布式ID生成器
func InitSnowFlake() {
	ip, err := GetOutBoundIP()
	if err != nil {
		zap.L().Error("Get Out Bound Ip error", zap.Error(err))
		panic(err)
	}

	var workId int64
	// 是否存在以本地IP为key的记录
	ifExist := model.GetIpExist(ip)
	if ifExist == 1 {
		// 如果存在则从redis获取对应的工作ID
		curWorkId, err := model.GetIp(ip)
		if err != nil {
			zap.L().Error("Get Ip Error", zap.Error(err))
			panic(err)
		}

		workId, err = strconv.ParseInt(curWorkId, 10, 64)
	} else {
		// 如果不存在，则生成新工作ID，并将其存储
		newWorkId, err := model.IncrWorkId(WorkID)
		if err != nil {
			zap.L().Error("New WorkId Error", zap.Error(err))
			panic(err)
		}
		// ip和工作ID映射
		model.SetWorkIdMapToIp(ip, newWorkId)
		workId = newWorkId
	}
	once.Do(func() {
		res, err := newSnowFlake(workId, 0)
		if err != nil {
			panic(err)
		}
		snowFlake = res
	})
}

func newSnowFlake(workId, datacenterId int64) (*Snowflake, error) {
	if workId < 0 || workId > maxWorkerId {
		return nil, errors.New("worker id out of range")
	}
	if datacenterId < 0 || datacenterId > maxDatacenterId {
		return nil, errors.New("datacenter id out of range ")
	}
	return &Snowflake{
		lastTimestamp: 0,
		workerId:      workId,
		datacenterId:  datacenterId,
		sequence:      0,
	}, nil
}

func NewSnowFlake() *Snowflake {
	if snowFlake == nil {
		once.Do(func() {
			res, _ := newSnowFlake(10, 10)
			snowFlake = res
		})
	}
	return snowFlake
}

func (sf *Snowflake) NextId() (int64, error) {
	sf.mu.Lock()
	defer sf.mu.Unlock()

	timestamp := time.Now().UnixNano() / 1000000

	if timestamp < sf.lastTimestamp {
		return 0, errors.New("clock moved backwards")
	}

	if timestamp == sf.lastTimestamp {
		sf.sequence = (sf.sequence + 1) & sequenceMask
		if sf.sequence == 0 {
			// 时钟回拨
			for timestamp <= sf.lastTimestamp {
				timestamp = time.Now().UnixNano() / 1000000
			}
		}
	} else {
		sf.sequence = 0
	}

	sf.lastTimestamp = timestamp
	// 相当于
	id := ((timestamp - twepoch) << timestampLeftShift) | (sf.datacenterId << datacenterIdShift) | (sf.workerId << workerIdShift) | sf.sequence

	return id, nil
}
