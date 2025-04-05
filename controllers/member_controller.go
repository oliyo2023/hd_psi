package controllers

import (
	"hd_psi/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type MemberController struct {
	db *gorm.DB
}

func NewMemberController(db *gorm.DB) *MemberController {
	return &MemberController{db: db}
}

func (mc *MemberController) ListMembers(c *gin.Context) {
	var members []models.Member
	if err := mc.db.Find(&members).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, members)
}

func (mc *MemberController) GetMember(c *gin.Context) {
	id := c.Param("id")
	var member models.Member
	if err := mc.db.First(&member, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Member not found"})
		return
	}
	c.JSON(http.StatusOK, member)
}

func (mc *MemberController) CreateMember(c *gin.Context) {
	var member models.Member
	if err := c.ShouldBindJSON(&member); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := mc.db.Create(&member).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, member)
}

func (mc *MemberController) UpdateMember(c *gin.Context) {
	id := c.Param("id")
	var member models.Member
	if err := mc.db.First(&member, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Member not found"})
		return
	}

	if err := c.ShouldBindJSON(&member); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := mc.db.Save(&member).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, member)
}

func (mc *MemberController) DeleteMember(c *gin.Context) {
	id := c.Param("id")
	if err := mc.db.Delete(&models.Member{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Member deleted"})
}
