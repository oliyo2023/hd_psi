package controllers

import (
	"hd_psi/models"
	"hd_psi/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AuthController struct {
	db *gorm.DB
}

func NewAuthController(db *gorm.DB) *AuthController {
	return &AuthController{db: db}
}

// LoginInput 定义用户登录请求的参数结构
type LoginInput struct {
	// Username 用户登录名，必填字段
	Username string `json:"username" binding:"required"`
	// Password 用户密码，必填字段
	Password string `json:"password" binding:"required"`
}

// RegisterInput 定义用户注册请求的参数结构
type RegisterInput struct {
	// Username 用户登录名，必填字段，必须唯一
	Username string `json:"username" binding:"required"`
	// Password 用户密码，必填字段
	Password string `json:"password" binding:"required"`
	// Name 用户真实姓名，必填字段
	Name string `json:"name" binding:"required"`
	// Email 用户电子邮箱，可选字段
	Email string `json:"email"`
	// Phone 用户手机号码，可选字段
	Phone string `json:"phone"`
	// Role 用户角色，必填字段
	Role models.Role `json:"role" binding:"required"`
	// StoreID 用户所属店铺ID，可选字段，总部人员可为空
	StoreID *uint `json:"store_id"`
}

// LoginResponse 定义登录成功后的响应数据结构
type LoginResponse struct {
	// Token JWT认证令牌，用于后续请求的身份验证
	Token string `json:"token"`
	// User 登录用户的基本信息，不包含敏感数据如密码
	User models.User `json:"user"`
	// ExpiresAt 令牌过期时间
	ExpiresAt time.Time `json:"expires_at"`
}

// Login 处理用户登录请求
// 验证用户凭证并生成JWT令牌
// 参数：
//   - c: Gin上下文对象，包含请求和响应信息
func (ac *AuthController) Login(c *gin.Context) {
	// 解析请求参数
	var input LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 查找用户
	var user models.User
	if err := ac.db.Where("username = ?", input.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
		return
	}

	// 检查用户状态
	if !user.Status {
		c.JSON(http.StatusForbidden, gin.H{"error": "用户已被禁用"})
		return
	}

	// 验证密码
	if !user.CheckPassword(input.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
		return
	}

	// 生成JWT令牌
	token, err := utils.GenerateToken(user.ID, user.Username, string(user.Role))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "生成令牌失败"})
		return
	}

	// 更新最后登录时间
	now := time.Now()
	user.LastLogin = &now
	ac.db.Save(&user)

	// 清除密码
	user.Password = ""

	// 返回令牌和用户信息
	c.JSON(http.StatusOK, LoginResponse{
		Token:     token,
		User:      user,
		ExpiresAt: time.Now().Add(24 * time.Hour),
	})
}

// Register 处理用户注册请求
// 创建新用户并将其保存到数据库
// 参数：
//   - c: Gin上下文对象，包含请求和响应信息
func (ac *AuthController) Register(c *gin.Context) {
	// 解析请求参数
	var input RegisterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 检查用户名是否已存在
	var existingUser models.User
	if err := ac.db.Where("username = ?", input.Username).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "用户名已存在"})
		return
	}

	// 创建新用户
	user := models.User{
		Username: input.Username,
		Password: input.Password, // 密码会在BeforeSave钩子中自动加密
		Name:     input.Name,
		Email:    input.Email,
		Phone:    input.Phone,
		Role:     input.Role,
		StoreID:  input.StoreID,
		Status:   true,
	}

	if err := ac.db.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建用户失败"})
		return
	}

	// 清除密码
	user.Password = ""

	// 返回创建成功的用户信息
	c.JSON(http.StatusCreated, user)
}

// GetProfile 获取当前登录用户的个人信息
// 从请求上下文中获取用户ID并返回用户详细信息
// 参数：
//   - c: Gin上下文对象，包含请求和响应信息
func (ac *AuthController) GetProfile(c *gin.Context) {
	// 从上下文中获取用户ID
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	// 查询用户信息
	var user models.User
	if err := ac.db.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		return
	}

	// 清除密码
	user.Password = ""

	// 返回用户信息
	c.JSON(http.StatusOK, user)
}

// UpdateProfile 更新当前登录用户的个人信息
// 允许用户修改姓名、邮箱和电话等基本信息
// 参数：
//   - c: Gin上下文对象，包含请求和响应信息
func (ac *AuthController) UpdateProfile(c *gin.Context) {
	// 从上下文中获取用户ID
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	// 查询用户信息
	var user models.User
	if err := ac.db.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		return
	}

	// 绑定请求数据
	var input struct {
		// Name 用户姓名
		Name string `json:"name"`
		// Email 电子邮箱
		Email string `json:"email"`
		// Phone 手机号码
		Phone string `json:"phone"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 更新用户信息
	user.Name = input.Name
	user.Email = input.Email
	user.Phone = input.Phone

	if err := ac.db.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新用户信息失败"})
		return
	}

	// 清除密码
	user.Password = ""

	// 返回更新后的用户信息
	c.JSON(http.StatusOK, user)
}

// ChangePassword 处理用户修改密码请求
// 验证旧密码并更新为新密码
// 参数：
//   - c: Gin上下文对象，包含请求和响应信息
func (ac *AuthController) ChangePassword(c *gin.Context) {
	// 从上下文中获取用户ID
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	// 查询用户信息
	var user models.User
	if err := ac.db.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		return
	}

	// 绑定请求数据
	var input struct {
		// OldPassword 用户当前密码，必填字段
		OldPassword string `json:"old_password" binding:"required"`
		// NewPassword 用户新密码，必填字段
		NewPassword string `json:"new_password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 验证旧密码
	if !user.CheckPassword(input.OldPassword) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "旧密码错误"})
		return
	}

	// 更新密码
	user.Password = input.NewPassword

	if err := ac.db.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新密码失败"})
		return
	}

	// 返回成功消息
	c.JSON(http.StatusOK, gin.H{"message": "密码修改成功"})
}
