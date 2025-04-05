package main

import (
	"hd_psi/config"
	"hd_psi/controllers"
	"hd_psi/models"
	"hd_psi/routes"

	"log"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	// 初始化数据库连接
	dsn := config.GetDBConfig()
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("数据库连接失败: ", err)
	}

	// 自动迁移数据模型
	db.AutoMigrate(
		&models.Product{},
		&models.Inventory{},
		&models.Supplier{},
		&models.PurchaseOrder{},
		&models.Store{},
		&models.InventoryTransaction{},
		&models.InventoryAlert{},
		&models.InventoryThreshold{},
		&models.Member{},
		&models.InventoryCheck{},
		&models.InventoryCheckItem{},
		&models.InventoryCheckAdjustment{},
		&models.SalesOrder{},
		&models.SalesOrderItem{},
		&models.NegotiationRecord{},
		&models.FittingRecord{},
		&models.ReturnOrder{},
		&models.ReturnOrderItem{},
		&models.FittingRoom{},
		&controllers.PointsTransaction{},
	)
	// 初始化Gin引擎
	r := gin.Default()

	// 注册路由
	routes.RegisterRoutes(r, db)

	// 启动服务
	if err := r.Run(":8080"); err != nil {
		log.Fatal("服务启动失败: ", err)
	}
}
