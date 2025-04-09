package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"hd_psi/backend/config"
	"hd_psi/backend/models"
	"io"
	"log"
	"os"
	"strconv"
	"strings"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	csvFile string
)

func main() {
	// 解析命令行参数
	flag.StringVar(&csvFile, "file", "", "CSV文件路径 (必填)")
	flag.Parse()

	// 验证必填参数
	if csvFile == "" {
		fmt.Println("错误: 请提供CSV文件路径")
		printUsage()
		os.Exit(1)
	}

	// 打开CSV文件
	file, err := os.Open(csvFile)
	if err != nil {
		log.Fatalf("无法打开CSV文件: %v", err)
	}
	defer file.Close()

	// 连接数据库
	dsn := config.GetDBConfig()
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("数据库连接失败: %v", err)
	}

	// 读取CSV文件
	reader := csv.NewReader(file)
	reader.Comma = ',' // 设置分隔符
	reader.LazyQuotes = true

	// 读取标题行
	header, err := reader.Read()
	if err != nil {
		log.Fatalf("读取CSV标题行失败: %v", err)
	}

	// 验证CSV格式
	requiredColumns := []string{"sku", "name"}
	for _, col := range requiredColumns {
		found := false
		for _, h := range header {
			if strings.ToLower(h) == col {
				found = true
				break
			}
		}
		if !found {
			log.Fatalf("CSV文件缺少必要的列: %s", col)
		}
	}

	// 获取列索引
	colIndex := make(map[string]int)
	for i, h := range header {
		colIndex[strings.ToLower(h)] = i
	}

	// 读取并处理每一行
	lineNum := 1 // 标题行是第1行
	successCount := 0
	errorCount := 0

	for {
		lineNum++
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Printf("警告: 第%d行读取失败: %v\n", lineNum, err)
			errorCount++
			continue
		}

		// 提取数据
		sku := getColumnValue(record, colIndex, "sku")
		name := getColumnValue(record, colIndex, "name")

		// 验证必填字段
		if sku == "" || name == "" {
			fmt.Printf("警告: 第%d行缺少必填字段 (SKU或名称)\n", lineNum)
			errorCount++
			continue
		}

		// 检查SKU是否已存在
		var existingProduct models.Product
		if err := db.Where("sku = ?", sku).First(&existingProduct).Error; err == nil {
			fmt.Printf("警告: 第%d行的SKU '%s' 已存在\n", lineNum, sku)
			errorCount++
			continue
		}

		// 创建商品
		product := models.Product{
			SKU:      sku,
			Name:     name,
			Color:    getColumnValue(record, colIndex, "color"),
			Size:     getColumnValue(record, colIndex, "size"),
			Season:   getColumnValue(record, colIndex, "season"),
			Category: getColumnValue(record, colIndex, "category"),
			Image:    getColumnValue(record, colIndex, "image"),
		}

		// 处理价格
		if costStr := getColumnValue(record, colIndex, "cost"); costStr != "" {
			if cost, err := strconv.ParseFloat(costStr, 64); err == nil {
				product.CostPrice = cost
			} else {
				fmt.Printf("警告: 第%d行的成本价格格式无效: %s\n", lineNum, costStr)
			}
		}

		if retailStr := getColumnValue(record, colIndex, "retail"); retailStr != "" {
			if retail, err := strconv.ParseFloat(retailStr, 64); err == nil {
				product.RetailPrice = retail
			} else {
				fmt.Printf("警告: 第%d行的零售价格格式无效: %s\n", lineNum, retailStr)
			}
		}

		// 保存到数据库
		if err := db.Create(&product).Error; err != nil {
			fmt.Printf("错误: 第%d行保存失败: %v\n", lineNum, err)
			errorCount++
			continue
		}

		fmt.Printf("成功: 第%d行 - 商品 '%s' (SKU: %s) 已添加\n", lineNum, product.Name, product.SKU)
		successCount++
	}

	// 打印总结
	fmt.Printf("\n导入完成: 成功 %d 条, 失败 %d 条\n", successCount, errorCount)
}

// 获取列值，不区分大小写
func getColumnValue(record []string, colIndex map[string]int, colName string) string {
	if idx, ok := colIndex[colName]; ok && idx < len(record) {
		return strings.TrimSpace(record[idx])
	}
	return ""
}

// 打印使用说明
func printUsage() {
	fmt.Println("\n使用方法:")
	fmt.Println("  batch_add_products -file <csv文件路径>")
	fmt.Println("\nCSV文件格式:")
	fmt.Println("  必须包含标题行，且必须包含'sku'和'name'列")
	fmt.Println("  支持的列: sku, name, color, size, season, category, image, cost, retail")
	fmt.Println("\n示例:")
	fmt.Println("  batch_add_products -file products.csv")
}
