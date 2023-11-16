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

func CreateUser(context *gin.Context){
	var input AuthenticateInput
	err := context.ShouldBindJSON(&input)

	if checkError(context, err){
        return
	}

	user := database.User{
        Username: input.Username,
        Password: input.Password,
    }
	savedUser, err := user.Create()

	if (checkError(context, err)){
		return
	}

	context.JSON(http.StatusCreated, gin.H{"user": savedUser})
}

func Authenticate(context *gin.Context){
	var input AuthenticateInput
	err := context.ShouldBindJSON(&input)

	if checkError(context, err){
        return
	}

	user, err := database.FindUserByUsername(input.Username)

	if checkError(context, err){
        return
	}

	err = user.ValidatePassword(input.Password)

	if checkError(context, err){
        return
	}

	jwt, err := utils.GenerateJWT(user)

	if checkError(context, err){
        return
	}

	context.JSON(http.StatusOK, gin.H{"jwt": jwt})
}

func checkError(context *gin.Context, err error) bool{
	isErr := err != nil
	if isErr{
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	return isErr
}