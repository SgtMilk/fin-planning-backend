package database

import (
	"errors"
	"regexp"
	"strings"

	"gorm.io/gorm"
)

type MonthlyExpense struct {
	gorm.Model

	Title        string  `gorm:"size:255;not null" json:"title"`
	StartMonth   string  `gorm:"size:7;not null" json:"startMonth"`
	EndMonth     string  `gorm:"size:7;not null" json:"endMonth"`
	CurrentValue float64 `gorm:"not null" json:"currentValue"`
	CPM          float64 `gorm:"not null" json:"cpm"`
	CIPY         float32 `gorm:"not null" json:"cipy"`
	IPY          float32 `gorm:"not null" json:"ipy"`
	IsTaxed      bool    `gorm:"not null" json:"isTaxed"`
	IsEditable   bool    `gorm:"not null" json:"isEditable"`
	Category     string  `gorm:"size:255;not null" json:"category"`

	// connections
	// to connect back from has many
	UserID uint `gorm:"not null" json:"userID"`

	// has many relationship
	Expenses []Expense `json:"expenses"`
}

func GetAllMonthlyExpense(userId uint) ([]MonthlyExpense, error) {
	var monthlyExpense []MonthlyExpense
	err := Database.Where("user_id=?", userId).Find(&monthlyExpense).Error
	return monthlyExpense, err
}

func GetMonthlyExpense(userId uint, monthlyExpenseId uint) (*MonthlyExpense, error) {
	var monthlyExpense MonthlyExpense
	err := Database.Where("id=? AND user_id=?", monthlyExpenseId, userId).Find(&monthlyExpense).Error
	return &monthlyExpense, err
}

func CreateMonthlyExpense(userID uint, category string) (*MonthlyExpense, error) {
	monthlyExpense := &MonthlyExpense{
		Title:        "New Value",
		StartMonth:   GetCurrentMonth(0),
		EndMonth:     GetCurrentMonth(0),
		CurrentValue: 0,
		CPM:          0,
		CIPY:         4,
		IPY:          0,
		IsTaxed:      false,
		IsEditable:   true,
		Category:     category,
		UserID:       userID,
	}

	err := Database.Create(&monthlyExpense).Error

	return monthlyExpense, err
}

func CreateMonthlyExpenses(monthlyExpenses []MonthlyExpense) error {
	return Database.Create(&monthlyExpenses).Error
}

func SaveMonthlyExpenses(monthlyExpenses []MonthlyExpense) error {
	return Database.Save(&monthlyExpenses).Error
}

func (monthlyExpense *MonthlyExpense) Delete() error {
	var expenses []Expense
	err := Database.Where("monthly_expense_id = ?", monthlyExpense.ID).Find(&expenses).Error
	if err != nil {
		return err
	}

	for _, expense := range expenses {
		err = expense.Delete()
		if err != nil {
			return err
		}
	}

	err = Database.Delete(&monthlyExpense).Error
	return err
}

func (monthlyExpense *MonthlyExpense) Update() error {
	// checking if the monthly expense is owned by the user
	_, err := GetMonthlyExpense(monthlyExpense.UserID, monthlyExpense.ID)
	if err != nil {
		return err
	}

	Database.Save(&monthlyExpense)
	return nil
}

func (monthlyExpense *MonthlyExpense) AddExpenses(expenses []Expense) error {
	for _, expense := range expenses {
		expense.MonthlyExpenseID = monthlyExpense.ID
	}

	err := CreateExpenses(expenses)
	return err
}

// ==============================================
// 					HOOKS
// ==============================================

func (monthlyExpense *MonthlyExpense) BeforeSave(*gorm.DB) error {
	return monthlyExpense.AssertInput()
}

// ==============================================
// 					HELPERS
// ==============================================

func (monthlyExpense *MonthlyExpense) AssertInput() error {
	// check if string lengths are right
	for _, str := range [2]string{monthlyExpense.Title, monthlyExpense.Category} {
		if len(str) > 64 {
			return errors.New("input string is too long")
		}
	}

	// check if dates are valid
	r, err := regexp.Compile("^[0-9]{4}-[0-9]{2}$")
	if err != nil {
		return err
	}

	for _, date := range [2]string{monthlyExpense.StartMonth, monthlyExpense.EndMonth} {
		matched := r.MatchString(date)
		if !matched {
			return errors.New("date has wrong format")
		}
	}

	if strings.Compare(monthlyExpense.StartMonth, monthlyExpense.EndMonth) < 0 {
		return errors.New("end month comes before start month")
	}

	return nil
}
