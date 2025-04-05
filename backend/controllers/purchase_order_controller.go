package controllers

import (
	"hd_psi/backend/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type PurchaseOrderController struct {
	db *gorm.DB
}

func NewPurchaseOrderController(db *gorm.DB) *PurchaseOrderController {
	return &PurchaseOrderController{db: db}
}

func (poc *PurchaseOrderController) ListPurchaseOrders(c *gin.Context) {
	var purchaseOrders []models.PurchaseOrder
	if err := poc.db.Find(&purchaseOrders).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, purchaseOrders)
}

func (poc *PurchaseOrderController) GetPurchaseOrder(c *gin.Context) {
	id := c.Param("id")
	var purchaseOrder models.PurchaseOrder
	if err := poc.db.First(&purchaseOrder, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "PurchaseOrder not found"})
		return
	}
	c.JSON(http.StatusOK, purchaseOrder)
}

func (poc *PurchaseOrderController) CreatePurchaseOrder(c *gin.Context) {
	var purchaseOrder models.PurchaseOrder
	if err := c.ShouldBindJSON(&purchaseOrder); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := poc.db.Create(&purchaseOrder).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, purchaseOrder)
}

func (poc *PurchaseOrderController) UpdatePurchaseOrder(c *gin.Context) {
	id := c.Param("id")
	var purchaseOrder models.PurchaseOrder
	if err := poc.db.First(&purchaseOrder, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "PurchaseOrder not found"})
		return
	}

	if err := c.ShouldBindJSON(&purchaseOrder); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := poc.db.Save(&purchaseOrder).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, purchaseOrder)
}

func (poc *PurchaseOrderController) DeletePurchaseOrder(c *gin.Context) {
	id := c.Param("id")
	if err := poc.db.Delete(&models.PurchaseOrder{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "PurchaseOrder deleted"})
}
