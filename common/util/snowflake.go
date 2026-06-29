package util

import (
	"github.com/sony/sonyflake"
	"log"
	"time"
)

// 单例模式
// 防止 同一机器同一进程 先后初始化实例，内部序列化每次从零开始，导致得到id相同
var sf *sonyflake.Sonyflake

func init() {
	// 创建 sonyflake 实例
	//settings := sonyflake.Settings{
	//	StartTime: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC), // 自定义起始时间
	//	MachineID: func() (uint16, error) {
	//		// 自定义机器ID获取逻辑，比如从环境变量或配置文件读取
	//		return 1, nil
	//	},
	//}
	// 可以自定义起始时间，这里使用默认（2014-09-01 00:00:00 +0000 UTC）
	sf = sonyflake.NewSonyflake(sonyflake.Settings{})
}

// 防止时间回拨
var lastTime int64 = 0

// GenerateSnowflake 雪花算法
func GenerateSnowflake() int64 {

	id, err := sf.NextID()
	if err != nil {
		log.Fatalf("生成ID失败: %v", err)
	}

	// 简单检测时间回拨
	now := time.Now().UnixNano() / 1e6 // 毫秒
	if now < lastTime {
		log.Printf("警告: 检测到时间回拨, lastTime=%d, now=%d", lastTime, now)
	}
	lastTime = now

	return int64(id)

}
