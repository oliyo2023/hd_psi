package controllers

import (
	"hd_psi/backend/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type FittingController struct {
	db *gorm.DB
}

func NewFittingController(db *gorm.DB) *FittingController {
	return &FittingController{db: db}
}

// ListFittingRooms 获取试衣间列表
func (fc *FittingController) ListFittingRooms(c *gin.Context) {
	var rooms []models.FittingRoom
	
	// 获取查询参数
	storeID := c.Query("store_id")
	status := c.Query("status")
	
	// 构建查询
	query := fc.db.Model(&models.FittingRoom{})
	
	if storeID != "" {
		query = query.Where("store_id = ?", storeID)
	}
	
	if status != "" {
		query = query.Where("status = ?", status)
	}
	
	// 执行查询
	if err := query.Find(&rooms).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, rooms)
}

// GetFittingRoom 获取试衣间详情
func (fc *FittingController) GetFittingRoom(c *gin.Context) {
	id := c.Param("id")
	var room models.FittingRoom
	if err := fc.db.First(&room, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Fitting room not found"})
		return
	}
	
	c.JSON(http.StatusOK, room)
}

// CreateFittingRoom 创建试衣间
func (fc *FittingController) CreateFittingRoom(c *gin.Context) {
	var room models.FittingRoom
	if err := c.ShouldBindJSON(&room); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	if err := fc.db.Create(&room).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusCreated, room)
}

// UpdateFittingRoom 更新试衣间
func (fc *FittingController) UpdateFittingRoom(c *gin.Context) {
	id := c.Param("id")
	var room models.FittingRoom
	if err := fc.db.First(&room, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Fitting room not found"})
		return
	}
	
	if err := c.ShouldBindJSON(&room); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	if err := fc.db.Save(&room).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, room)
}

// DeleteFittingRoom 删除试衣间
func (fc *FittingController) DeleteFittingRoom(c *gin.Context) {
	id := c.Param("id")
	if err := fc.db.Delete(&models.FittingRoom{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "Fitting room deleted"})
}

// ListFittingRecords 获取试衣记录列表
func (fc *FittingController) ListFittingRecords(c *gin.Context) {
	var records []models.FittingRecord
	
	// 获取查询参数
	memberID := c.Query("member_id")
	storeID := c.Query("store_id")
	productID := c.Query("product_id")
	
	// 构建查询
	query := fc.db.Model(&models.FittingRecord{})
	
	if memberID != "" {
		query = query.Where("member_id = ?", memberID)
	}
	
	if storeID != "" {
		query = query.Where("store_id = ?", storeID)
	}
	
	if productID != "" {
		query = query.Where("product_id = ?", productID)
	}
	
	// 执行查询
	if err := query.Order("created_at DESC").Find(&records).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, records)
}

// GetFittingRecord 获取试衣记录详情
func (fc *FittingController) GetFittingRecord(c *gin.Context) {
	id := c.Param("id")
	var record models.FittingRecord
	if err := fc.db.First(&record, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Fitting record not found"})
		return
	}
	
	c.JSON(http.StatusOK, record)
}

// CreateFittingRecord 创建试衣记录
func (fc *FittingController) CreateFittingRecord(c *gin.Context) {
	var record models.FittingRecord
	if err := c.ShouldBindJSON(&record); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	// 验证会员是否存在
	var member models.Member
	if err := fc.db.First(&member, record.MemberID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Member not found"})
		return
	}
	
	// 验证商品是否存在
	var product models.Product
	if err := fc.db.First(&product, record.ProductID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}
	
	// 验证试衣间是否存在
	var room models.FittingRoom
	if err := fc.db.First(&room, record.FittingRoomID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Fitting room not found"})
		return
	}
	
	// 验证店铺是否存在
	var store models.Store
	if err := fc.db.First(&store, record.StoreID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Store not found"})
		return
	}
	
	// 创建试衣记录
	if err := fc.db.Create(&record).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	// 更新试衣间状态为占用
	room.Status = "occupied"
	if err := fc.db.Save(&room).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update fitting room status"})
		return
	}
	
	c.JSON(http.StatusCreated, record)
}

// UpdateFittingRecord 更新试衣记录
func (fc *FittingController) UpdateFittingRecord(c *gin.Context) {
	id := c.Param("id")
	var record models.FittingRecord
	if err := fc.db.First(&record, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Fitting record not found"})
		return
	}
	
	// 保存原试衣间ID
	oldRoomID := record.FittingRoomID
	
	if err := c.ShouldBindJSON(&record); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	// 开始事务
	tx := fc.db.Begin()
	
	// 如果试衣间ID变更，更新试衣间状态
	if oldRoomID != record.FittingRoomID {
		// 将原试衣间状态改为可用
		var oldRoom models.FittingRoom
		if err := tx.First(&oldRoom, oldRoomID).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get old fitting room"})
			return
		}
		
		oldRoom.Status = "available"
		if err := tx.Save(&oldRoom).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update old fitting room status"})
			return
		}
		
		// 将新试衣间状态改为占用
		var newRoom models.FittingRoom
		if err := tx.First(&newRoom, record.FittingRoomID).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get new fitting room"})
			return
		}
		
		newRoom.Status = "occupied"
		if err := tx.Save(&newRoom).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update new fitting room status"})
			return
		}
	}
	
	// 更新试衣记录
	if err := tx.Save(&record).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	// 提交事务
	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
		return
	}
	
	c.JSON(http.StatusOK, record)
}

// CompleteFitting 完成试衣
func (fc *FittingController) CompleteFitting(c *gin.Context) {
	id := c.Param("id")
	var record models.FittingRecord
	if err := fc.db.First(&record, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Fitting record not found"})
		return
	}
	
	var input struct {
		SatisfactionLevel int    `json:"satisfaction_level"`
		Comments          string `json:"comments"`
	}
	
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	// 开始事务
	tx := fc.db.Begin()
	
	// 更新试衣记录
	record.SatisfactionLevel = input.SatisfactionLevel
	record.Comments = input.Comments
	
	if err := tx.Save(&record).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	// 将试衣间状态改为可用
	var room models.FittingRoom
	if err := tx.First(&room, record.FittingRoomID).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get fitting room"})
		return
	}
	
	room.Status = "available"
	if err := tx.Save(&room).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update fitting room status"})
		return
	}
	
	// 提交事务
	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
		return
	}
	
	c.JSON(http.StatusOK, record)
}
