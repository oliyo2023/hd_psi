package controllers

import (
	"hd_psi/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type InventoryAlertController struct {
	db *gorm.DB
}

func NewInventoryAlertController(db *gorm.DB) *InventoryAlertController {
	return &InventoryAlertController{db: db}
}

// ListAlerts 获取所有预警
func (iac *InventoryAlertController) ListAlerts(c *gin.Context) {
	var alerts []models.InventoryAlert
	
	// 获取查询参数
	status := c.Query("status")
	storeID := c.Query("store_id")
	alertType := c.Query("alert_type")
	
	// 构建查询
	query := iac.db.Model(&models.InventoryAlert{})
	
	if status != "" {
		query = query.Where("status = ?", status)
	}
	
	if storeID != "" {
		query = query.Where("store_id = ?", storeID)
	}
	
	if alertType != "" {
		query = query.Where("alert_type = ?", alertType)
	}
	
	// 执行查询
	if err := query.Find(&alerts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, alerts)
}

// GetAlert 获取单个预警
func (iac *InventoryAlertController) GetAlert(c *gin.Context) {
	id := c.Param("id")
	var alert models.InventoryAlert
	if err := iac.db.First(&alert, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Alert not found"})
		return
	}
	c.JSON(http.StatusOK, alert)
}

// UpdateAlertStatus 更新预警状态
func (iac *InventoryAlertController) UpdateAlertStatus(c *gin.Context) {
	id := c.Param("id")
	var alert models.InventoryAlert
	if err := iac.db.First(&alert, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Alert not found"})
		return
	}
	
	// 绑定请求数据
	var input struct {
		Status string `json:"status" binding:"required"`
	}
	
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	// 更新状态
	alert.Status = models.AlertStatus(input.Status)
	
	// 如果状态为已解决，设置解决时间
	if alert.Status == models.Resolved {
		now := time.Now()
		alert.ResolvedAt = &now
	}
	
	if err := iac.db.Save(&alert).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, alert)
}

// CheckInventoryLevels 检查库存水平并生成预警
// 这个方法可以通过定时任务调用，或者在库存变动后调用
func (iac *InventoryAlertController) CheckInventoryLevels(c *gin.Context) {
	// 获取所有库存记录
	var inventories []struct {
		StoreID   uint
		ProductID uint
		Category  string
		Quantity  int
	}
	
	// 联合查询获取库存和商品类别
	if err := iac.db.Table("inventories").
		Select("inventories.store_id, inventories.product_id, products.category, inventories.quantity").
		Joins("JOIN products ON inventories.product_id = products.id").
		Scan(&inventories).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	// 获取所有预警阈值设置
	var thresholds []struct {
		StoreID   uint
		Category  string
		LowLevel  int
		HighLevel int
	}
	
	// 这里假设有一个阈值设置表，实际中需要创建这个表
	// 如果没有这个表，可以使用默认值或配置文件中的设置
	if err := iac.db.Table("inventory_thresholds").Find(&thresholds).Error; err != nil {
		// 如果表不存在，使用默认值
		thresholds = []struct {
			StoreID   uint
			Category  string
			LowLevel  int
			HighLevel int
		}{
			{StoreID: 0, Category: "", LowLevel: 10, HighLevel: 100}, // 默认阈值
		}
	}
	
	// 检查每个库存记录
	var newAlerts []models.InventoryAlert
	for _, inv := range inventories {
		// 查找适用的阈值
		var threshold struct {
			LowLevel  int
			HighLevel int
		}
		
		// 首先尝试找到特定店铺和类别的阈值
		found := false
		for _, t := range thresholds {
			if (t.StoreID == inv.StoreID || t.StoreID == 0) && 
			   (t.Category == inv.Category || t.Category == "") {
				threshold.LowLevel = t.LowLevel
				threshold.HighLevel = t.HighLevel
				found = true
				break
			}
		}
		
		if !found {
			// 使用默认阈值
			threshold.LowLevel = 10
			threshold.HighLevel = 100
		}
		
		// 检查是否低于最低阈值
		if inv.Quantity <= threshold.LowLevel {
			// 检查是否已存在活跃的低库存预警
			var existingAlert models.InventoryAlert
			result := iac.db.Where("store_id = ? AND product_id = ? AND alert_type = ? AND status = ?", 
				inv.StoreID, inv.ProductID, models.LowStock, models.Active).First(&existingAlert)
			
			if result.Error == gorm.ErrRecordNotFound {
				// 创建新预警
				newAlerts = append(newAlerts, models.InventoryAlert{
					StoreID:     inv.StoreID,
					ProductID:   inv.ProductID,
					Category:    inv.Category,
					AlertType:   models.LowStock,
					Threshold:   threshold.LowLevel,
					CurrentQty:  inv.Quantity,
					Status:      models.Active,
					Description: "库存低于最低阈值",
				})
			}
		} else if inv.Quantity >= threshold.HighLevel {
			// 检查是否已存在活跃的高库存预警
			var existingAlert models.InventoryAlert
			result := iac.db.Where("store_id = ? AND product_id = ? AND alert_type = ? AND status = ?", 
				inv.StoreID, inv.ProductID, models.Overstock, models.Active).First(&existingAlert)
			
			if result.Error == gorm.ErrRecordNotFound {
				// 创建新预警
				newAlerts = append(newAlerts, models.InventoryAlert{
					StoreID:     inv.StoreID,
					ProductID:   inv.ProductID,
					Category:    inv.Category,
					AlertType:   models.Overstock,
					Threshold:   threshold.HighLevel,
					CurrentQty:  inv.Quantity,
					Status:      models.Active,
					Description: "库存超过最高阈值",
				})
			}
		}
	}
	
	// 批量创建新预警
	if len(newAlerts) > 0 {
		if err := iac.db.Create(&newAlerts).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}
	
	c.JSON(http.StatusOK, gin.H{
		"message": "库存检查完成",
		"new_alerts": len(newAlerts),
	})
}
