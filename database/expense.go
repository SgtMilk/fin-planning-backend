package database

import (
	"time"

	"gorm.io/gorm"
)

type Expense struct{
	gorm.Model
	
	Title string `gorm:"size:255;not null" json:"title"`
	Amount float64 `gorm:"not null" json:"amount"`
	Date time.Time `gorm:"not null" json:"date"`

	// connections
	// to connect back from has many
	MonthlyExpenseID uint `gorm:"not null" json:"monthlyExpenseID"`
}