package models

import (
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Role 定义系统中的用户角色类型
type Role string

// 系统中的用户角色常量定义
const (
	// Admin 管理员角色，拥有系统的全部权限
	Admin Role = "admin"
	// Manager 店长角色，负责管理单个店铺的运营
	Manager Role = "manager"
	// Staff 普通员工角色，执行日常业务操作
	Staff Role = "staff"
	// Cashier 收银员角色，负责处理销售和收款
	Cashier Role = "cashier"
	// Operator 操作员角色，负责特定的系统操作任务
	Operator Role = "operator"
)

// User 定义系统用户模型，存储用户账户信息和权限
type User struct {
	// ID 用户唯一标识
	ID uint `gorm:"primaryKey"`
	// Username 用户登录名，必须唯一
	Username string `gorm:"size:50;uniqueIndex;not null"`
	// Password 用户密码，存储加密后的哈希值
	Password string `gorm:"size:255;not null"`
	// Name 用户真实姓名
	Name string `gorm:"size:50;not null"`
	// Email 用户电子邮箱
	Email string `gorm:"size:100"`
	// Phone 用户手机号码
	Phone string `gorm:"size:20"`
	// Role 用户角色，决定用户权限
	Role Role `gorm:"size:20;not null;default:'staff'"`
	// StoreID 用户所属店铺ID，可以为空表示总部人员
	StoreID *uint `gorm:"default:null"`
	// LastLogin 用户最后登录时间
	LastLogin *time.Time
	// Status 用户状态，true表示启用，false表示禁用
	Status bool `gorm:"default:true"`
	// CreatedAt 用户创建时间
	CreatedAt time.Time
	// UpdatedAt 用户信息最后更新时间
	UpdatedAt time.Time
}

// BeforeSave 在用户信息保存到数据库前自动加密密码
// 这是 GORM 的钩子函数，在创建或更新用户时自动调用
// 参数：
//   - tx: GORM 数据库事务对象
//
// 返回：
//   - error: 如果加密过程出错则返回错误，否则返回 nil
func (u *User) BeforeSave(tx *gorm.DB) error {
	if len(u.Password) > 0 {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		u.Password = string(hashedPassword)
	}
	return nil
}

// CheckPassword 验证用户密码是否正确
// 将用户输入的原始密码与数据库中存储的加密密码进行比对
// 参数：
//   - password: 用户输入的原始密码
//
// 返回：
//   - bool: 如果密码正确返回 true，否则返回 false
func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}
