package database

import "gorm.io/gorm"

type MonthlyExpense struct{
	gorm.Model

	Title string `gorm:"size:255;not null" json:"title"`
	StartMonth string `gorm:"size:7;not null" json:"startMonth"`
	EndMonth string `gorm:"size:7;not null" json:"endMonth"`
	CurrentValue float64 `gorm:"not null" json:"currentValue"`
	CPM float64 `gorm:"not null" json:"cpm"`
	CIPY float32 `gorm:"not null" json:"cipy"`
	IPY float32 `gorm:"not null" json:"ipy"`
	IsTaxed bool `gorm:"not null" json:"isTaxed"`
	IsEditable bool `gorm:"not null" json:"isEditable"`
	Category string `gorm:"size:255;not null" json:"category"`

	// connections
	// to connect back from has many
	UserID uint `gorm:"not null" json:"userID"`

	// has many relationship
	Expenses []Expense `json:"expenses"`
}

func (monthlyExpense *MonthlyExpense) Delete() error{
	var expenses []Expense
	err := Database.Where("monthly_expense_id = ?", monthlyExpense.ID).Find(&expenses).Error
	if err != nil{
		return err
	}

	for _, expense := range expenses {
		err = expense.Delete()
		if err != nil{
			return err
		}
	}

	err = Database.Delete(&monthlyExpense).Error
	return err
}