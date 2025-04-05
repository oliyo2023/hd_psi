package models

type Member struct {
	ID    uint `gorm:"primaryKey"`
	Name  string
	Phone string
}
