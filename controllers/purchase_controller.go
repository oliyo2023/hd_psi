package controllers

import (
	"hd_psi/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type PurchaseController struct {
	db *gorm.DB
}

func NewPurchaseController(db *gorm.DB) *PurchaseController {
	return &PurchaseController{db: db}
}

func (pc *PurchaseController) ListPurchaseOrders(c *gin.Context) {
	var purchaseOrders []models.PurchaseOrder
	if err := pc.db.Find(&purchaseOrders).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, purchaseOrders)
}

func (pc *PurchaseController) GetPurchaseOrder(c *gin.Context) {
	id := c.Param("id")
	var purchaseOrder models.PurchaseOrder
	if err := pc.db.First(&purchaseOrder, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "PurchaseOrder not found"})
		return
	}
	c.JSON(http.StatusOK, purchaseOrder)
}

func (pc *PurchaseController) CreatePurchaseOrder(c *gin.Context) {
	var purchaseOrder models.PurchaseOrder
	if err := c.ShouldBindJSON(&purchaseOrder); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := pc.db.Create(&purchaseOrder).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, purchaseOrder)
}

func (pc *PurchaseController) UpdatePurchaseOrder(c *gin.Context) {
	id := c.Param("id")
	var purchaseOrder models.PurchaseOrder
	if err := pc.db.First(&purchaseOrder, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "PurchaseOrder not found"})
		return
	}

	if err := c.ShouldBindJSON(&purchaseOrder); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := pc.db.Save(&purchaseOrder).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, purchaseOrder)
}

func (pc *PurchaseController) DeletePurchaseOrder(c *gin.Context) {
	id := c.Param("id")
	if err := pc.db.Delete(&models.PurchaseOrder{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "PurchaseOrder deleted"})
}
