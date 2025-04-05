package controllers

import (
	"hd_psi/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// PointsTransaction 积分交易记录
type PointsTransaction struct {
	ID          uint      `gorm:"primaryKey"`
	MemberID    uint      `gorm:"not null"` // 会员ID
	Points      int       `gorm:"not null"` // 积分变动，正数为增加，负数为减少
	Type        string    `gorm:"size:20;not null"` // 类型：购物/兑换/调整
	ReferenceID *uint     // 关联单据ID
	ReferenceType string  `gorm:"size:20"` // 关联单据类型
	Description string    `gorm:"size:255"` // 描述
	OperatorID  uint      // 操作人ID
	CreatedAt   time.Time
}

type MemberPointsController struct {
	db *gorm.DB
}

func NewMemberPointsController(db *gorm.DB) *MemberPointsController {
	return &MemberPointsController{db: db}
}

// GetMemberPoints 获取会员积分
func (mpc *MemberPointsController) GetMemberPoints(c *gin.Context) {
	memberID := c.Param("id")
	var member models.Member
	if err := mpc.db.First(&member, memberID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Member not found"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"member_id": member.ID,
		"name": member.Name,
		"points": member.Points,
	})
}

// ListPointsTransactions 获取会员积分交易记录
func (mpc *MemberPointsController) ListPointsTransactions(c *gin.Context) {
	memberID := c.Param("id")
	
	// 验证会员是否存在
	var member models.Member
	if err := mpc.db.First(&member, memberID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Member not found"})
		return
	}
	
	var transactions []PointsTransaction
	if err := mpc.db.Where("member_id = ?", memberID).Order("created_at DESC").Find(&transactions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, transactions)
}

// AddPoints 增加会员积分
func (mpc *MemberPointsController) AddPoints(c *gin.Context) {
	memberID := c.Param("id")
	
	var input struct {
		Points       int    `json:"points" binding:"required"`
		Type         string `json:"type" binding:"required"`
		ReferenceID  *uint  `json:"reference_id"`
		ReferenceType string `json:"reference_type"`
		Description  string `json:"description"`
		OperatorID   uint   `json:"operator_id" binding:"required"`
	}
	
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	// 验证积分必须为正数
	if input.Points <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Points must be positive"})
		return
	}
	
	// 开始事务
	tx := mpc.db.Begin()
	
	// 获取会员信息
	var member models.Member
	if err := tx.First(&member, memberID).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusNotFound, gin.H{"error": "Member not found"})
		return
	}
	
	// 更新会员积分
	member.Points += input.Points
	if err := tx.Save(&member).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update member points"})
		return
	}
	
	// 创建积分交易记录
	transaction := PointsTransaction{
		MemberID:      member.ID,
		Points:        input.Points,
		Type:          input.Type,
		ReferenceID:   input.ReferenceID,
		ReferenceType: input.ReferenceType,
		Description:   input.Description,
		OperatorID:    input.OperatorID,
		CreatedAt:     time.Now(),
	}
	
	if err := tx.Create(&transaction).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create points transaction"})
		return
	}
	
	// 提交事务
	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"member_id": member.ID,
		"name": member.Name,
		"points": member.Points,
		"transaction": transaction,
	})
}

// DeductPoints 扣减会员积分
func (mpc *MemberPointsController) DeductPoints(c *gin.Context) {
	memberID := c.Param("id")
	
	var input struct {
		Points       int    `json:"points" binding:"required"`
		Type         string `json:"type" binding:"required"`
		ReferenceID  *uint  `json:"reference_id"`
		ReferenceType string `json:"reference_type"`
		Description  string `json:"description"`
		OperatorID   uint   `json:"operator_id" binding:"required"`
	}
	
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	// 验证积分必须为正数
	if input.Points <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Points must be positive"})
		return
	}
	
	// 开始事务
	tx := mpc.db.Begin()
	
	// 获取会员信息
	var member models.Member
	if err := tx.First(&member, memberID).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusNotFound, gin.H{"error": "Member not found"})
		return
	}
	
	// 检查积分是否足够
	if member.Points < input.Points {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{"error": "Insufficient points"})
		return
	}
	
	// 更新会员积分
	member.Points -= input.Points
	if err := tx.Save(&member).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update member points"})
		return
	}
	
	// 创建积分交易记录
	transaction := PointsTransaction{
		MemberID:      member.ID,
		Points:        -input.Points, // 负数表示扣减
		Type:          input.Type,
		ReferenceID:   input.ReferenceID,
		ReferenceType: input.ReferenceType,
		Description:   input.Description,
		OperatorID:    input.OperatorID,
		CreatedAt:     time.Now(),
	}
	
	if err := tx.Create(&transaction).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create points transaction"})
		return
	}
	
	// 提交事务
	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"member_id": member.ID,
		"name": member.Name,
		"points": member.Points,
		"transaction": transaction,
	})
}

// CalculateMemberLevel 计算会员等级
func (mpc *MemberPointsController) CalculateMemberLevel(c *gin.Context) {
	memberID := c.Param("id")
	
	// 获取会员信息
	var member models.Member
	if err := mpc.db.First(&member, memberID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Member not found"})
		return
	}
	
	// 根据累计消费金额计算会员等级
	var newLevel models.MemberLevel
	
	switch {
	case member.TotalSpent >= 50000:
		newLevel = models.Diamond
	case member.TotalSpent >= 20000:
		newLevel = models.Platinum
	case member.TotalSpent >= 10000:
		newLevel = models.Gold
	case member.TotalSpent >= 5000:
		newLevel = models.Silver
	default:
		newLevel = models.Regular
	}
	
	// 如果等级有变化，更新会员等级
	if member.Level != newLevel {
		member.Level = newLevel
		if err := mpc.db.Save(&member).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update member level"})
			return
		}
	}
	
	c.JSON(http.StatusOK, gin.H{
		"member_id": member.ID,
		"name": member.Name,
		"level": member.Level,
		"total_spent": member.TotalSpent,
	})
}
