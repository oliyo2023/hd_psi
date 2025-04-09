package main

import (
	"flag"
	"fmt"
	"hd_psi/backend/config"
	"hd_psi/backend/models"
	"log"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// 定义命令行参数
var (
	sku         string
	name        string
	color       string
	size        string
	season      string
	category    string
	image       string
	costPrice   float64
	retailPrice float64
)

func main() {
	// 解析命令行参数
	flag.StringVar(&sku, "sku", "", "商品SKU编码 (必填)")
	flag.StringVar(&name, "name", "", "商品名称 (必填)")
	flag.StringVar(&color, "color", "", "商品颜色")
	flag.StringVar(&size, "size", "", "商品尺码")
	flag.StringVar(&season, "season", "", "商品季节")
	flag.StringVar(&category, "category", "", "商品类别")
	flag.StringVar(&image, "image", "", "商品图片URL")
	flag.Float64Var(&costPrice, "cost", 0, "成本价")
	flag.Float64Var(&retailPrice, "retail", 0, "零售价")

	flag.Parse()

	// 验证必填参数
	if sku == "" || name == "" {
		fmt.Println("错误: SKU和商品名称为必填项")
		printUsage()
		os.Exit(1)
	}

	// 连接数据库
	dsn := config.GetDBConfig()
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("数据库连接失败: %v", err)
	}

	// 检查SKU是否已存在
	var existingProduct models.Product
	if err := db.Where("sku = ?", sku).First(&existingProduct).Error; err == nil {
		fmt.Printf("错误: SKU '%s' 已存在\n", sku)
		os.Exit(1)
	}

	// 创建商品
	product := models.Product{
		SKU:         sku,
		Name:        name,
		Color:       color,
		Size:        size,
		Season:      season,
		Category:    category,
		Image:       image,
		CostPrice:   costPrice,
		RetailPrice: retailPrice,
	}

	if err := db.Create(&product).Error; err != nil {
		log.Fatalf("创建商品失败: %v", err)
	}

	fmt.Printf("商品创建成功! ID: %d, SKU: %s, 名称: %s\n", product.ID, product.SKU, product.Name)
}

// 打印使用说明
func printUsage() {
	fmt.Println("\n使用方法:")
	fmt.Println("  create_product [选项]")
	fmt.Println("\n选项:")
	fmt.Println("  -sku string      商品SKU编码 (必填)")
	fmt.Println("  -name string     商品名称 (必填)")
	fmt.Println("  -color string    商品颜色")
	fmt.Println("  -size string     商品尺码")
	fmt.Println("  -season string   商品季节")
	fmt.Println("  -category string 商品类别")
	fmt.Println("  -image string    商品图片URL")
	fmt.Println("  -cost float      成本价")
	fmt.Println("  -retail float    零售价")
	fmt.Println("\n示例:")
	fmt.Printf("  %s -sku \"MS001\" -name \"男士休闲衬衫\" -color \"蓝色\" -size \"XL\" -category \"衬衫\" -cost 80 -retail 199\n", os.Args[0])
}
