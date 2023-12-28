package database

import (
	"strconv"
	"time"
)

func GetCurrentMonth(monthOffset int) string{
	year, month, _ := time.Now().AddDate(0, monthOffset, 0).Date()
	curMonth := int(month)
	curMonthString := strconv.Itoa(curMonth)
	if curMonth < 10{
		curMonthString = "0" + curMonthString
	}

	return strconv.Itoa(year) + "-" + curMonthString	
}