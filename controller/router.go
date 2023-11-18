package controller

import (
	"github.com/SgtMilk/fin-planning-backend/middleware"
	"github.com/gin-gonic/gin"
)

func CreateRouter(noLogger bool) *gin.Engine{
    var router *gin.Engine
    if noLogger{
        router = gin.New()
    } else{
        router = gin.Default()
    }

    publicRoutes := router.Group("/auth")
    publicRoutes.POST("/register", CreateUser)
    publicRoutes.POST("/login", Authenticate)

    protectedRoutes := router.Group("/api")
    protectedRoutes.Use(middleware.JWTAuthMiddleware())
    publicRoutes.DELETE("/closeaccount", Authenticate)

    return router
}