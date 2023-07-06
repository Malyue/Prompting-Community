package snowflake

//type Snowflake struct {
//	mu            sync.Mutex
//	lastTimestamp int64
//	workerId      int64
//	datacenterId  int64
//	sequence      int64
//}
//
//var (
//	snowFlake *Snowflake
//	once      sync.Once
//)
//
//const (
//	twepoch            = int64(1417937700000) // Unix纪元时间戳
//	workerIdBits       = uint(5)              // 机器ID所占位数
//	datacenterBits     = uint(5)              // 数据中心ID所占位数
//	maxWorkerId        = int64(-1) ^ (int64(-1) << workerIdBits)
//	maxDatacenterId    = int64(-1) ^ (int64(-1) << datacenterBits)
//	sequenceBits       = uint(12) // 序列号所占位数
//	workerIdShift      = sequenceBits
//	datacenterIdShift  = sequenceBits + workerIdBits
//	timestampLeftShift = sequenceBits + workerIdBits + datacenterBits
//	sequenceMask       = int64(-1) ^ (int64(-1) << sequenceBits)
//)
//
//func InitSnowFlake() {
//	ip, err := network.GetOutBoundIP()
//	if err != nil {
//		zap.L().Error("Get Out Bound Ip error", zap.Error(err))
//		panic(err)
//	}
//
//	var workId int64
//	ctx := context.Background()
//
//}
