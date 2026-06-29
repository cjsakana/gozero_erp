package excel

import (
	"fmt"
	"github.com/xuri/excelize/v2"
	"google.golang.org/protobuf/types/known/timestamppb"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// FieldMapping 定义字段映射关系
type FieldMapping struct {
	ExcelColumn string // Excel 列名
	StructField string // 结构体字段名
}

// ParseOptions 解析选项
type ParseOptions struct {
	FieldMappings []FieldMapping // 自定义字段映射
	IgnoreCase    bool           // 是否忽略大小写
}

// ParseExcelToStruct 将 Excel 文件解析成结构体切片
// filePath: Excel 文件路径
// sheetName: 工作表名称
// headerRow: 表头所在行（从 1 开始）
// targetStruct: 目标结构体实例（用于反射）
// options: 解析选项（可选）
// return: 结构体切片，错误信息
func ParseExcelToStruct(filePath, sheetName string, headerRow int, targetStruct any, options ...*ParseOptions) ([]any, error) {
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open Excel file: %v", err)
	}
	defer f.Close()

	// 如果没有指定表名，则获取第一个工作表名
	if sheetName == "" {
		sheets := f.GetSheetList()
		if len(sheets) == 0 {
			return nil, fmt.Errorf("excel is empty, sheets list is empty")
		}
		sheetName = sheets[0]
	}

	// 获取所有行
	rows, err := f.GetRows(sheetName)
	if err != nil {
		return nil, fmt.Errorf("failed to get rows: %v", err)
	}

	if len(rows) < headerRow {
		return nil, fmt.Errorf("header row %d not found", headerRow)
	}

	// 获取表头（列名）
	headers := rows[headerRow-1]

	// 获取目标结构体的反射信息
	targetType := reflect.TypeOf(targetStruct)
	if targetType.Kind() == reflect.Ptr {
		targetType = targetType.Elem()
	}

	// 检查是否是结构体
	if targetType.Kind() != reflect.Struct {
		return nil, fmt.Errorf("targetStruct must be a struct or pointer to struct")
	}

	// 处理选项
	var opt *ParseOptions
	if len(options) > 0 {
		opt = options[0]
	} else {
		opt = &ParseOptions{}
	}

	// 创建字段映射关系
	fieldMap := make(map[string]int) // Excel列名 -> 结构体字段索引

	for i := 0; i < targetType.NumField(); i++ {
		field := targetType.Field(i)

		// 1. 首先检查自定义映射
		var excelColName string
		if opt.FieldMappings != nil {
			for _, mapping := range opt.FieldMappings {
				if mapping.StructField == field.Name {
					excelColName = mapping.ExcelColumn
					break
				}
			}
		}

		// 2. 检查结构体标签
		if excelColName == "" {
			if tag := field.Tag.Get("excel"); tag != "" {
				excelColName = tag
			}
		}

		// 3. 默认使用字段名
		if excelColName == "" {
			excelColName = field.Name
		}

		// 记录映射关系
		fieldMap[excelColName] = i
	}

	// 存储解析后的数据
	var results []any

	// 从数据行开始遍历（跳过表头）
	for i := headerRow; i < len(rows); i++ {
		row := rows[i]

		// 创建新的结构体实例
		newStruct := reflect.New(targetType).Elem()

		// 遍历Excel列
		for colIndex, header := range headers {
			if colIndex >= len(row) {
				continue
			}

			cellValue := row[colIndex]
			headerName := header

			// 查找对应的结构体字段
			var fieldIndex int
			var found bool

			// 尝试精确匹配
			if idx, ok := fieldMap[headerName]; ok {
				fieldIndex = idx
				found = true
			}

			// 尝试忽略大小写匹配
			if !found && opt.IgnoreCase {
				for excelName, idx := range fieldMap {
					if strings.EqualFold(excelName, headerName) {
						fieldIndex = idx
						found = true
						break
					}
				}
			}

			if !found {
				continue
			}

			field := targetType.Field(fieldIndex)

			// 根据字段类型设置值
			switch field.Type.Kind() {
			case reflect.String:
				newStruct.Field(fieldIndex).SetString(cellValue)

			case reflect.Int32, reflect.Int64:
				if num, ok := parseUnixFromCell(cellValue); ok {
					if field.Type.Kind() == reflect.Int32 {
						newStruct.Field(fieldIndex).SetInt(int64(int32(num)))
					} else {
						newStruct.Field(fieldIndex).SetInt(num)
					}
				}

			case reflect.Float32, reflect.Float64:
				if num, err := strconv.ParseFloat(cellValue, 64); err == nil {
					newStruct.Field(fieldIndex).SetFloat(num)
				}

			case reflect.Bool:
				if b, err := parseBoolValue(cellValue); err == nil {
					newStruct.Field(fieldIndex).SetBool(b)
				} else {
					// 可以记录错误或使用默认值
					newStruct.Field(fieldIndex).SetBool(false)
					// 或者返回错误
					// return fmt.Errorf("字段 %s: %v", field.Name, err)
				}

			case reflect.Struct:
				// 处理 protobuf 生成的时间类型
				if field.Type == reflect.TypeOf(timestamppb.Timestamp{}) {
					if cellValue != "" {
						timeFloat, err := strconv.ParseFloat(cellValue, 64)
						if err == nil {
							timeTime, err := excelize.ExcelDateToTime(timeFloat, false)
							if err == nil {
								ts := timestamppb.New(timeTime)
								newStruct.Field(fieldIndex).Set(reflect.ValueOf(ts))
							}
						}
					}
				}
				// 可以添加其他 protobuf 特殊类型的处理

			default:
				// 其他类型可以在这里处理或忽略
			}
		}

		// 添加到结果
		results = append(results, newStruct.Interface())
	}

	return results, nil
}

func parseUnixFromCell(cellValue string) (int64, bool) {
	cellValue = strings.TrimSpace(cellValue)
	if cellValue == "" {
		return 0, false
	}

	if num, err := strconv.ParseInt(cellValue, 10, 64); err == nil {
		return num, true
	}

	if f, err := strconv.ParseFloat(cellValue, 64); err == nil {
		if t, err2 := excelize.ExcelDateToTime(f, false); err2 == nil {
			return t.Unix(), true
		}
	}

	layouts := []string{
		"2006/1/2", "2006/01/02", "2006-1-2", "2006-01-02", "2006.1.2", "2006.01.02",
		"2006/1/2 15:04:05", "2006/01/02 15:04:05", "2006-1-2 15:04:05", "2006-01-02 15:04:05",
		"2006.1.2 15:04:05", "2006.01.02 15:04:05",
	}
	for _, layout := range layouts {
		if t, err := time.ParseInLocation(layout, cellValue, time.Local); err == nil {
			return t.Unix(), true
		}
	}
	return 0, false
}

func parseBoolValue(cellValue string) (bool, error) {
	// 去除前后空格
	cellValue = strings.TrimSpace(cellValue)
	// 空字符串视为false
	if cellValue == "" {
		return false, nil
	}

	// 检查常见的true表示
	if matchesTrue(cellValue) {
		return true, nil
	}

	// 检查常见的false表示
	if matchesFalse(cellValue) {
		return false, nil
	}

	// 尝试解析为数字
	if num, err := strconv.Atoi(cellValue); err == nil {
		return num != 0, nil
	}

	// 无法识别的格式
	return false, fmt.Errorf("无法识别的布尔值: %s", cellValue)
}

func matchesTrue(value string) bool {
	trueValues := []string{"true", "t", "yes", "y", "是", "1", "真"}
	value = strings.ToLower(value)
	for _, v := range trueValues {
		if value == v {
			return true
		}
	}
	return false
}

func matchesFalse(value string) bool {
	falseValues := []string{"false", "f", "no", "n", "否", "0", "假"}
	value = strings.ToLower(value)
	for _, v := range falseValues {
		if value == v {
			return true
		}
	}
	return false
}
