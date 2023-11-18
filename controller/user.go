package controller

import (
	"errors"
	"net/http"

	"github.com/SgtMilk/fin-planning-backend/database"
	"github.com/SgtMilk/fin-planning-backend/utils"
	"github.com/gin-gonic/gin"
	passwordvalidator "github.com/wagslane/go-password-validator"
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

	err = assertInput(&input)

	if checkError(context, err){
        return
	}

	user := database.User{
        Username: input.Username,
        Password: input.Password,
    }
	_ , err = user.Create()

	if checkError(context, err){
		return
	}

	context.JSON(http.StatusCreated, nil)
}

func DeleteUser(context *gin.Context){
	user, err := utils.GetCurrentUser(context)

	if checkError(context, err){
		return
	}

	err = user.Delete()

	if checkError(context, err){
		return
	}

	context.JSON(http.StatusOK, nil)
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

func assertInput(input *AuthenticateInput) error{
	// username 
	usernameLengthCheck := len(input.Username) < 8 || len(input.Username) > 256
	passwordLengthCheck := len(input.Password) < 8 || len(input.Password) > 72
	if usernameLengthCheck && passwordLengthCheck{
		return errors.New("username and password not of right size")
	}else if usernameLengthCheck{
		return errors.New("username not of right size")
	}else if passwordLengthCheck{
		return errors.New("password not of right size")
	}

	// strength evaluation
	err := passwordvalidator.Validate(input.Username, 50)
	if err != nil{
		return errors.New("insecure username, try using a longer username")
	}

	return passwordvalidator.Validate(input.Password, 60)
}