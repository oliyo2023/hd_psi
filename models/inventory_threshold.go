package models

import "time"

// InventoryThreshold 库存阈值设置
type InventoryThreshold struct {
	ID        uint   `gorm:"primaryKey"`
	StoreID   uint   `gorm:"default:0"` // 0表示适用于所有店铺
	Category  string `gorm:"size:50;default:''"` // 空字符串表示适用于所有类别
	LowLevel  int    `gorm:"not null"` // 低库存阈值
	HighLevel int    `gorm:"not null"` // 高库存阈值
	CreatedAt time.Time
	UpdatedAt time.Time
}
