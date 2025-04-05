package models

type Product struct {
	ID          uint   `gorm:"primaryKey"`
	SKU         string `gorm:"size:255;uniqueIndex"`
	Name        string
	Color       string
	Size        string
	Season      string
	Category    string
	Image       string
	CostPrice   float64
	RetailPrice float64
}
