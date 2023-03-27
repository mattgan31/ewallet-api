package models

type Transaction struct {
	GormModel
	Type        string `gorm:"not_null" json:"type" validate:"required"`
	Amount      int    `gorm:"not_null" json:"amount" validate:"required"`
	Description string `gorm:"not_null" json:"description" validate:"required"`
	User_ID     int
}
