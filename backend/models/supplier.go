package models

type Supplier struct {
	ID            uint `gorm:"primaryKey"`
	Name          string
	Qualification string
	Rating        string
	City          string
}
