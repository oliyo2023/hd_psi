package controllers

import (
	"hd_psi/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type InventoryThresholdController struct {
	db *gorm.DB
}

func NewInventoryThresholdController(db *gorm.DB) *InventoryThresholdController {
	return &InventoryThresholdController{db: db}
}

// ListThresholds 获取所有库存阈值设置
func (itc *InventoryThresholdController) ListThresholds(c *gin.Context) {
	var thresholds []models.InventoryThreshold
	
	// 获取查询参数
	storeID := c.Query("store_id")
	category := c.Query("category")
	
	// 构建查询
	query := itc.db.Model(&models.InventoryThreshold{})
	
	if storeID != "" {
		query = query.Where("store_id = ? OR store_id = 0", storeID)
	}
	
	if category != "" {
		query = query.Where("category = ? OR category = ''", category)
	}
	
	// 执行查询
	if err := query.Find(&thresholds).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, thresholds)
}

// GetThreshold 获取单个阈值设置
func (itc *InventoryThresholdController) GetThreshold(c *gin.Context) {
	id := c.Param("id")
	var threshold models.InventoryThreshold
	if err := itc.db.First(&threshold, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Threshold not found"})
		return
	}
	c.JSON(http.StatusOK, threshold)
}

// CreateThreshold 创建阈值设置
func (itc *InventoryThresholdController) CreateThreshold(c *gin.Context) {
	var threshold models.InventoryThreshold
	if err := c.ShouldBindJSON(&threshold); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	// 验证高阈值必须大于低阈值
	if threshold.HighLevel <= threshold.LowLevel {
		c.JSON(http.StatusBadRequest, gin.H{"error": "High level must be greater than low level"})
		return
	}
	
	if err := itc.db.Create(&threshold).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusCreated, threshold)
}

// UpdateThreshold 更新阈值设置
func (itc *InventoryThresholdController) UpdateThreshold(c *gin.Context) {
	id := c.Param("id")
	var threshold models.InventoryThreshold
	if err := itc.db.First(&threshold, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Threshold not found"})
		return
	}
	
	// 绑定请求数据
	if err := c.ShouldBindJSON(&threshold); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	// 验证高阈值必须大于低阈值
	if threshold.HighLevel <= threshold.LowLevel {
		c.JSON(http.StatusBadRequest, gin.H{"error": "High level must be greater than low level"})
		return
	}
	
	if err := itc.db.Save(&threshold).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, threshold)
}

// DeleteThreshold 删除阈值设置
func (itc *InventoryThresholdController) DeleteThreshold(c *gin.Context) {
	id := c.Param("id")
	if err := itc.db.Delete(&models.InventoryThreshold{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "Threshold deleted"})
}
