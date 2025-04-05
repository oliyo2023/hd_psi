package controllers

import (
	"hd_psi/backend/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type SupplierController struct {
	db *gorm.DB
}

func NewSupplierController(db *gorm.DB) *SupplierController {
	return &SupplierController{db: db}
}

// 供应商列表请求参数
type ListSuppliersQuery struct {
	Name   string `form:"name"`
	Code   string `form:"code"`
	Type   string `form:"type"`
	Status string `form:"status"`
	Page   int    `form:"page,default=1"`
	Limit  int    `form:"limit,default=10"`
}

// 供应商列表响应
type SuppliersResponse struct {
	Total int               `json:"total"`
	Items []models.Supplier `json:"items"`
}

// ListSuppliers 获取供应商列表
func (sc *SupplierController) ListSuppliers(c *gin.Context) {
	var query ListSuppliersQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 构建查询
	db := sc.db.Model(&models.Supplier{})

	// 应用过滤条件
	if query.Name != "" {
		db = db.Where("name LIKE ?", "%"+query.Name+"%")
	}
	if query.Code != "" {
		db = db.Where("code LIKE ?", "%"+query.Code+"%")
	}
	if query.Type != "" {
		db = db.Where("type = ?", query.Type)
	}
	if query.Status != "" {
		status := query.Status == "active"
		db = db.Where("status = ?", status)
	}

	// 计算总数
	var total int64
	db.Count(&total)

	// 分页
	offset := (query.Page - 1) * query.Limit
	var suppliers []models.Supplier

	if err := db.Offset(offset).Limit(query.Limit).
		Order("created_at DESC").
		Find(&suppliers).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, SuppliersResponse{
		Total: int(total),
		Items: suppliers,
	})
}

// GetSupplier 获取供应商详情
func (sc *SupplierController) GetSupplier(c *gin.Context) {
	id := c.Param("id")
	var supplier models.Supplier

	if err := sc.db.First(&supplier, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "供应商不存在"})
		return
	}

	c.JSON(http.StatusOK, supplier)
}

// 创建供应商请求
type CreateSupplierRequest struct {
	Name          string                `json:"name" binding:"required"`
	Code          string                `json:"code" binding:"required"`
	Type          models.SupplierType   `json:"type" binding:"required"`
	ContactPerson string                `json:"contact_person"`
	ContactPhone  string                `json:"contact_phone"`
	Email         string                `json:"email"`
	Address       string                `json:"address"`
	City          string                `json:"city"`
	Rating        models.SupplierRating `json:"rating"`
	Qualification string                `json:"qualification"`
	PaymentTerms  string                `json:"payment_terms"`
	DeliveryTerms string                `json:"delivery_terms"`
	Status        bool                  `json:"status"`
	Note          string                `json:"note"`
}

// CreateSupplier 创建供应商
func (sc *SupplierController) CreateSupplier(c *gin.Context) {
	var request CreateSupplierRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 检查供应商编码是否已存在
	var existingSupplier models.Supplier
	if err := sc.db.Where("code = ?", request.Code).First(&existingSupplier).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "供应商编码已存在"})
		return
	}

	// 创建供应商
	supplier := models.Supplier{
		Name:          request.Name,
		Code:          request.Code,
		Type:          request.Type,
		ContactPerson: request.ContactPerson,
		ContactPhone:  request.ContactPhone,
		Email:         request.Email,
		Address:       request.Address,
		City:          request.City,
		Rating:        request.Rating,
		Qualification: request.Qualification,
		PaymentTerms:  request.PaymentTerms,
		DeliveryTerms: request.DeliveryTerms,
		Status:        request.Status,
		Note:          request.Note,
	}

	if err := sc.db.Create(&supplier).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建供应商失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, supplier)
}

// UpdateSupplier 更新供应商
func (sc *SupplierController) UpdateSupplier(c *gin.Context) {
	id := c.Param("id")
	var supplier models.Supplier

	// 查询供应商
	if err := sc.db.First(&supplier, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "供应商不存在"})
		return
	}

	var request CreateSupplierRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 检查供应商编码是否已被其他供应商使用
	if request.Code != supplier.Code {
		var existingSupplier models.Supplier
		if err := sc.db.Where("code = ? AND id != ?", request.Code, id).First(&existingSupplier).Error; err == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "供应商编码已被其他供应商使用"})
			return
		}
	}

	// 更新供应商
	supplier.Name = request.Name
	supplier.Code = request.Code
	supplier.Type = request.Type
	supplier.ContactPerson = request.ContactPerson
	supplier.ContactPhone = request.ContactPhone
	supplier.Email = request.Email
	supplier.Address = request.Address
	supplier.City = request.City
	supplier.Rating = request.Rating
	supplier.Qualification = request.Qualification
	supplier.PaymentTerms = request.PaymentTerms
	supplier.DeliveryTerms = request.DeliveryTerms
	supplier.Status = request.Status
	supplier.Note = request.Note

	if err := sc.db.Save(&supplier).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新供应商失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, supplier)
}

// DeleteSupplier 删除供应商
func (sc *SupplierController) DeleteSupplier(c *gin.Context) {
	id := c.Param("id")

	// 检查供应商是否存在
	var supplier models.Supplier
	if err := sc.db.First(&supplier, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "供应商不存在"})
		return
	}

	// 检查供应商是否被采购单引用
	var count int64
	sc.db.Model(&models.PurchaseOrder{}).Where("supplier_id = ?", id).Count(&count)
	if count > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "供应商已被采购单引用，无法删除"})
		return
	}

	// 删除供应商
	if err := sc.db.Delete(&supplier).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除供应商失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "供应商删除成功"})
}
