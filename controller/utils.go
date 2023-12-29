package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func CheckError(context *gin.Context, err error) bool {
	isErr := err != nil
	if isErr {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	return isErr
}

func FetchId(key string, context *gin.Context) (uint, error) {
	idString := context.Param(key)
	id, err := (strconv.ParseUint(idString, 10, 64))
	return uint(id), err
}
