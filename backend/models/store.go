package models

import "time"

type Store struct {
	ID        uint      `gorm:"primaryKey"`
	Name      string    `gorm:"size:100;not null"`
	Address   string    `gorm:"size:255"`
	Type      string    `gorm:"size:50"` // 男装店/女装店
	Phone     string    `gorm:"size:20"`
	Manager   string    `gorm:"size:50"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
