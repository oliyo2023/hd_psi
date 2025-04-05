package models

type PurchaseOrder struct {
	ID         uint `gorm:"primaryKey"`
	SupplierID uint
	ProductID  uint
	Quantity   int
	Price      float64
	Status     string
}
