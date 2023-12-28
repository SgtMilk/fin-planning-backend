package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func CheckError(context *gin.Context, err error) bool{
	isErr := err != nil
	if isErr{
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	return isErr
}