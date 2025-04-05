package routes

import (
	"hd_psi/controllers"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// RegisterRoutes 注册所有路由
func RegisterRoutes(r *gin.Engine, db *gorm.DB) {
	// 商品管理路由
	productController := controllers.NewProductController(db)
	productGroup := r.Group("/products")
	{
		productGroup.GET("", productController.ListProducts)
		productGroup.GET("/:id", productController.GetProduct)
		productGroup.POST("", productController.CreateProduct)
		productGroup.PUT("/:id", productController.UpdateProduct)
		productGroup.DELETE("/:id", productController.DeleteProduct)
	}

	// 库存管理路由
	inventoryController := controllers.NewInventoryController(db)
	inventoryGroup := r.Group("/inventory")
	{
		inventoryGroup.GET("", inventoryController.ListInventories)
		inventoryGroup.GET("/:id", inventoryController.GetInventory)
		inventoryGroup.POST("", inventoryController.CreateInventory)
		inventoryGroup.PUT("/:id", inventoryController.UpdateInventory)
		inventoryGroup.DELETE("/:id", inventoryController.DeleteInventory)
	}

	// 采购管理路由
	purchaseController := controllers.NewPurchaseController(db)
	purchaseGroup := r.Group("/purchases")
	{
		purchaseGroup.GET("", purchaseController.ListPurchaseOrders)
		purchaseGroup.GET("/:id", purchaseController.GetPurchaseOrder)
		purchaseGroup.POST("", purchaseController.CreatePurchaseOrder)
		purchaseGroup.PUT("/:id", purchaseController.UpdatePurchaseOrder)
		purchaseGroup.DELETE("/:id", purchaseController.DeletePurchaseOrder)
	}

	// 会员管理路由
	memberController := controllers.NewMemberController(db)
	memberGroup := r.Group("/members")
	{
		memberGroup.GET("", memberController.ListMembers)
		memberGroup.GET("/:id", memberController.GetMember)
		memberGroup.POST("", memberController.CreateMember)
		memberGroup.PUT("/:id", memberController.UpdateMember)
		memberGroup.DELETE("/:id", memberController.DeleteMember)

		// 会员积分路由
		memberPointsController := controllers.NewMemberPointsController(db)
		memberGroup.GET("/:id/points", memberPointsController.GetMemberPoints)
		memberGroup.GET("/:id/points/transactions", memberPointsController.ListPointsTransactions)
		memberGroup.POST("/:id/points/add", memberPointsController.AddPoints)
		memberGroup.POST("/:id/points/deduct", memberPointsController.DeductPoints)
		memberGroup.POST("/:id/level/calculate", memberPointsController.CalculateMemberLevel)
	}

	// 店铺管理路由
	storeController := controllers.NewStoreController(db)
	storeGroup := r.Group("/stores")
	{
		storeGroup.GET("", storeController.ListStores)
		storeGroup.GET("/:id", storeController.GetStore)
		storeGroup.POST("", storeController.CreateStore)
		storeGroup.PUT("/:id", storeController.UpdateStore)
		storeGroup.DELETE("/:id", storeController.DeleteStore)
	}

	// 库存交易路由
	inventoryTransactionController := controllers.NewInventoryTransactionController(db)
	transactionGroup := r.Group("/inventory-transactions")
	{
		transactionGroup.GET("", inventoryTransactionController.ListTransactions)
		transactionGroup.GET("/:id", inventoryTransactionController.GetTransaction)
		transactionGroup.POST("", inventoryTransactionController.CreateTransaction)
		transactionGroup.GET("/store/:storeId", inventoryTransactionController.GetStoreTransactions)
		transactionGroup.GET("/product/:productId", inventoryTransactionController.GetProductTransactions)
	}

	// 库存预警路由
	inventoryAlertController := controllers.NewInventoryAlertController(db)
	alertGroup := r.Group("/inventory-alerts")
	{
		alertGroup.GET("", inventoryAlertController.ListAlerts)
		alertGroup.GET("/:id", inventoryAlertController.GetAlert)
		alertGroup.PUT("/:id/status", inventoryAlertController.UpdateAlertStatus)
		alertGroup.POST("/check", inventoryAlertController.CheckInventoryLevels)
	}

	// 库存阈值路由
	inventoryThresholdController := controllers.NewInventoryThresholdController(db)
	thresholdGroup := r.Group("/inventory-thresholds")
	{
		thresholdGroup.GET("", inventoryThresholdController.ListThresholds)
		thresholdGroup.GET("/:id", inventoryThresholdController.GetThreshold)
		thresholdGroup.POST("", inventoryThresholdController.CreateThreshold)
		thresholdGroup.PUT("/:id", inventoryThresholdController.UpdateThreshold)
		thresholdGroup.DELETE("/:id", inventoryThresholdController.DeleteThreshold)
	}

	// 库存盘点路由
	inventoryCheckController := controllers.NewInventoryCheckController(db)
	checkGroup := r.Group("/inventory-checks")
	{
		checkGroup.GET("", inventoryCheckController.ListChecks)
		checkGroup.GET("/:id", inventoryCheckController.GetCheck)
		checkGroup.POST("", inventoryCheckController.CreateCheck)
		checkGroup.PUT("/:id/start", inventoryCheckController.StartCheck)
		checkGroup.PUT("/:id/complete", inventoryCheckController.CompleteCheck)
		checkGroup.PUT("/:id/cancel", inventoryCheckController.CancelCheck)
		checkGroup.PUT("/:id/items/:itemId", inventoryCheckController.UpdateCheckItem)
		checkGroup.POST("/:id/adjustments", inventoryCheckController.CreateAdjustment)
		checkGroup.PUT("/adjustments/:adjustmentId/approve", inventoryCheckController.ApproveAdjustment)
	}

	// 销售管理路由
	salesController := controllers.NewSalesController(db)
	salesGroup := r.Group("/sales")
	{
		// 销售订单路由
		salesGroup.GET("/orders", salesController.ListOrders)
		salesGroup.GET("/orders/:id", salesController.GetOrder)
		salesGroup.POST("/orders", salesController.CreateOrder)
		salesGroup.PUT("/orders/:id/status", salesController.UpdateOrderStatus)

		// 退换货路由
		salesGroup.POST("/returns", salesController.CreateReturnOrder)
		salesGroup.PUT("/returns/:id/status", salesController.UpdateReturnOrderStatus)
	}

	// 试衣管理路由
	fittingController := controllers.NewFittingController(db)
	fittingGroup := r.Group("/fitting")
	{
		// 试衣间路由
		fittingGroup.GET("/rooms", fittingController.ListFittingRooms)
		fittingGroup.GET("/rooms/:id", fittingController.GetFittingRoom)
		fittingGroup.POST("/rooms", fittingController.CreateFittingRoom)
		fittingGroup.PUT("/rooms/:id", fittingController.UpdateFittingRoom)
		fittingGroup.DELETE("/rooms/:id", fittingController.DeleteFittingRoom)

		// 试衣记录路由
		fittingGroup.GET("/records", fittingController.ListFittingRecords)
		fittingGroup.GET("/records/:id", fittingController.GetFittingRecord)
		fittingGroup.POST("/records", fittingController.CreateFittingRecord)
		fittingGroup.PUT("/records/:id", fittingController.UpdateFittingRecord)
		fittingGroup.PUT("/records/:id/complete", fittingController.CompleteFitting)
	}
}
