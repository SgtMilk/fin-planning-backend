package database

import (
	"strconv"
	"time"

	"gorm.io/gorm"
)

type Options struct{
	gorm.Model
	
	StartMonth string `gorm:"size:7" json:"startDate"`
	EndMonth string `gorm:"size:7" json:"endDate"`
	MonthInterval uint16 `json:"monthInterval"`
	Inflation float32 `json:"inflation"`
	TaxRate float32 `json:"taxRate"`

	// connections
	// has many relationship
	Balance []BalanceAttribute `json:"balance"`
}

func CreateDefaultOptions() (*Options, error){
	year, month, _ := time.Now().Date()
	curMonth := int(month)
	curMonthString := strconv.Itoa(curMonth)
	if curMonth < 10{
		curMonthString = "0" + curMonthString
	}

	options := &Options{
		StartMonth: strconv.Itoa(year) + "-" + curMonthString,
		EndMonth: strconv.Itoa(year + 50) + "-" + curMonthString,
		MonthInterval: 12,
		Inflation: 4,
		TaxRate: 50,
	}
	err := Database.Create(&options).Error

	if err != nil{
		return &Options{}, err
	}
	return options, nil
}

func (options *Options) Delete() error{
	var balanceAttributes []BalanceAttribute
	err := Database.Where("options_id = ?", options.ID).Find(&balanceAttributes).Error
	if err != nil{
		return err
	}

	for _, balanceAttribute := range balanceAttributes {
		err = balanceAttribute.Delete()
		if err != nil{
			return err
		}
	}

	err = Database.Delete(&options).Error
	return err
}