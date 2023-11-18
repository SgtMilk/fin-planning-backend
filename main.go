package main

import (
	"fmt"

	"github.com/SgtMilk/fin-planning-backend/controller"
	"github.com/SgtMilk/fin-planning-backend/database"
)

func main(){
	database.LoadDB()
	
	router := controller.CreateRouter(false)
	router.Run(":8000")
    fmt.Println("Server running on port 8000")
}