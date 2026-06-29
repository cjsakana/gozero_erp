package excel

import (
	"fmt"
)

// 示例结构体
type Product struct {
	ID       int64   `excel:"产品ID"`
	Name     string  `excel:"产品名称"`
	Price    float64 `excel:"单价"`
	IsActive bool    `excel:"是否有效"`
}

func ExampleParseExcelToStruct() {
	// 示例1：基本用法（使用结构体标签映射）
	products, err := ParseExcelToStruct("products.xlsx", "Sheet1", 1, Product{})
	if err != nil {
		fmt.Printf("解析失败: %v\n", err)
		return
	}

	for _, p := range products {
		product := p.(Product)
		fmt.Printf("ID: %d, Name: %s, Price: %.2f\n",
			product.ID, product.Name, product.Price)
	}

	// 示例2：使用外部映射配置
	options := &ParseOptions{
		FieldMappings: []FieldMapping{
			{"产品ID", "ID"},
			{"产品名称", "Name"},
			{"价格", "Price"},
			{"是否有效", "IsActive"},
		},
	}

	products, err = ParseExcelToStruct("products_v2.xlsx", "数据", 1, Product{}, options)
	if err != nil {
		fmt.Printf("解析失败: %v\n", err)
		return
	}

	for _, p := range products {
		product := p.(Product)
		fmt.Printf("ID: %d, Name: %s, Price: %.2f\n",
			product.ID, product.Name, product.Price)
	}
}
