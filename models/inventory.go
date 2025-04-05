package models

type Inventory struct {
	ID        uint `gorm:"primaryKey"`
	StoreID   uint
	ProductID uint
	Quantity  int
}
