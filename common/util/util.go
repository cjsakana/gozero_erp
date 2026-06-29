package util

import (
	"context"
	"crypto/rand"
	"erp/common/xcode"
	"fmt"
	"math/big"
	"path/filepath"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

// RandomNumeric 生成指定 size 的随机数字字符串
func RandomNumeric(size int) string {
	if size <= 0 {
		panic("size must be > 0")
	}
	b := make([]byte, size)
	for i := 0; i < size; i++ {
		n, _ := rand.Int(rand.Reader, big.NewInt(10))
		b[i] = byte('0' + n.Int64())
	}
	return string(b)
}

// IsEmpty 判断值是否为空的辅助函数
func IsEmpty(value interface{}) bool {
	switch v := value.(type) {
	case int:
		return v == 0
	case float64:
		return v == 0
	case string:
		return v == ""
	case time.Time:
		return v.IsZero()
	default:
		return false
	}
}

// ExtractBirthdayFromID18 从18位身份证提取出生日期
func ExtractBirthdayFromID18(idCard string) (time.Time, error) {
	if len(idCard) != 18 {
		return time.Time{}, fmt.Errorf("身份证号码长度不正确")
	}

	// 提取年月日部分：第7-14位
	birthdayStr := idCard[6:14]

	// 解析为时间
	birthday, err := time.Parse("20060102", birthdayStr)
	if err != nil {
		return time.Time{}, fmt.Errorf("身份证出生日期格式错误: %v", err)
	}

	return birthday, nil
}

var seq int64 // 全局自增序列（原子操作，线程安全）
// GenerateNo 生成业务单号（例如：ORD20251029-123456789001）
// prefix: 业务前缀，例如 "ORD"、"DEL"、"PUR"
// 格式: {prefix}{日期}-时间戳后9位 + 3位自增序列
func GenerateNo(prefix string) string {
	now := time.Now()
	date := now.Format("20060102")            // 日期部分
	timestamp := now.UnixNano() / 1e6         // 毫秒时间戳
	seqNum := atomic.AddInt64(&seq, 1) % 1000 // 每毫秒最多1000个不同序号
	lastDigits := timestamp % 1000000000      // 时间戳后9位
	return fmt.Sprintf("%s%s-%09d%03d", prefix, date, lastDigits, seqNum)
}

// GetFileExtension 返回带点的后缀，如 ".jpg" 要么是 ""
func GetFileExtension(filename string) string {
	ext := filepath.Ext(filename)
	// filepath.Ext 返回带点的后缀，如 ".jpg" 要么是 ""
	return ext
}

// GetInt64FromCtx   从context中获取 UserId/EmployeeId（字符串转int64）
// gozero 的 ctx 存储数值是 float64，需要 string 才能避免float64精度丢失
func GetInt64FromCtx(ctx context.Context, key string) (int64, error) {
	value := ctx.Value(key)
	if value == nil {
		return 0, xcode.New(401, "用户ID不存在")
	}

	// Id作为字符串存储在JWT中
	str, ok := value.(string)
	if !ok {
		return 0, xcode.New(401, "用户ID类型错误")
	}

	id, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return 0, xcode.New(401, "用户ID格式错误")
	}

	return id, nil
}

// ==================== ID转换工具函数 ====================

// StringToInt64 将string类型的ID转换为int64
// 如果转换失败，返回0和error
func StringToInt64(idStr string) (int64, error) {
	if idStr == "" || idStr == "0" {
		return 0, nil
	}
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		logx.Errorf("failed to parse id string to int64, idStr: %s, err: %v", idStr, err)
		return 0, err
	}
	return id, nil
}

// Int64ToString 将int64类型的ID转换为string
func Int64ToString(id int64) string {
	if id == 0 {
		return "0"
	}
	return strconv.FormatInt(id, 10)
}

// StringSliceToInt64Slice 将string切片转换为int64切片
func StringSliceToInt64Slice(strSlice []string) ([]int64, error) {
	if len(strSlice) == 0 {
		return []int64{}, nil
	}

	result := make([]int64, 0, len(strSlice))
	for _, str := range strSlice {
		id, err := StringToInt64(str)
		if err != nil {
			return nil, err
		}
		result = append(result, id)
	}
	return result, nil
}

// Int64SliceToStringSlice 将int64切片转换为string切片
func Int64SliceToStringSlice(int64Slice []int64) []string {
	if len(int64Slice) == 0 {
		return []string{}
	}

	result := make([]string, 0, len(int64Slice))
	for _, id := range int64Slice {
		result = append(result, Int64ToString(id))
	}
	return result
}
