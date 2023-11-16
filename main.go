package main

import (
	"github.com/SgtMilk/fin-planning-backend/controller"
	"github.com/SgtMilk/fin-planning-backend/database"
)

func main(){
	database.LoadDB()
	controller.ServeApplication()
}