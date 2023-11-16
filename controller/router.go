package controller

import (
	"fmt"

	"github.com/SgtMilk/fin-planning-backend/middleware"
	"github.com/gin-gonic/gin"
)

func ServeApplication() {
    router := gin.Default()

    publicRoutes := router.Group("/auth")
    publicRoutes.POST("/createuser", CreateUser)
    publicRoutes.POST("/authenticate", Authenticate)

    protectedRoutes := router.Group("/api")
    protectedRoutes.Use(middleware.JWTAuthMiddleware())

    router.Run(":8000")
    fmt.Println("Server running on port 8000")
}