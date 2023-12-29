package controller

import (
	"net/http"

	"github.com/SgtMilk/fin-planning-backend/database"
	"github.com/SgtMilk/fin-planning-backend/utils"
	"github.com/gin-gonic/gin"
)

type AuthenticateInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UpdatePasswordInput struct {
	OldPassword string `json:"oldPassword" binding:"required"`
	NewPassword string `json:"newPassword" binding:"required"`
}

func CreateUserRouter(router *gin.RouterGroup) {
	userRoutes := router.Group("/user")
	userRoutes.DELETE("/close_account", DeleteUser)
	userRoutes.PUT("/update_password", UpdatePassword)
}

func CreateUser(context *gin.Context) {
	var input AuthenticateInput
	err := context.ShouldBindJSON(&input)
	if CheckError(context, err) {
		return
	}

	user := database.User{
		Username: input.Username,
		Password: input.Password,
	}
	_, err = user.Create()
	if CheckError(context, err) {
		return
	}

	context.JSON(http.StatusCreated, nil)
}

func DeleteUser(context *gin.Context) {
	user, err := utils.GetCurrentUser(context)
	if CheckError(context, err) {
		return
	}

	err = user.Delete()
	if CheckError(context, err) {
		return
	}

	context.JSON(http.StatusOK, nil)
}

func UpdatePassword(context *gin.Context) {
	var input UpdatePasswordInput
	err := context.ShouldBindJSON(&input)
	if CheckError(context, err) {
		return
	}

	user, err := utils.GetCurrentUser(context)
	if CheckError(context, err) {
		return
	}

	err = user.UpdatePassword(input.OldPassword, input.NewPassword)
	if CheckError(context, err) {
		return
	}

	context.JSON(http.StatusOK, nil)
}

func Authenticate(context *gin.Context) {
	var input AuthenticateInput
	err := context.ShouldBindJSON(&input)
	if CheckError(context, err) {
		return
	}

	user, err := database.FindUserByUsername(input.Username)
	if CheckError(context, err) {
		return
	}

	err = user.ValidatePassword(input.Password)
	if CheckError(context, err) {
		return
	}

	jwt, err := utils.GenerateJWT(user)
	if CheckError(context, err) {
		return
	}

	context.JSON(http.StatusOK, gin.H{"jwt": jwt})
}
