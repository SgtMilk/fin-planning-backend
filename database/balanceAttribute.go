package database

import "gorm.io/gorm"

type BalanceAttribute struct{
	gorm.Model
	
	Proportion float32 `gorm:"not null" json:"proportion"`
	IsPositive bool `gorm:"not null" json:"isPositive"`

	// connections
	// to connect back from has many
	OptionsID uint `gorm:"not null" json:"optionsID"`
}

func (balanceAttribute *BalanceAttribute) Delete() error{
	err := Database.Delete(&balanceAttribute).Error
	return err
}