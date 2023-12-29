package controller

import (
	"errors"
	"net/http"

	"github.com/SgtMilk/fin-planning-backend/database"
	"github.com/SgtMilk/fin-planning-backend/utils"
	"github.com/gin-gonic/gin"
)

func CreateMonthlyExpenseRouter(router *gin.RouterGroup) {
	monthlyExpenseRoutes := router.Group("/monthly_expense")
	monthlyExpenseRoutes.POST("/", CreateMonthlyExpense)
	monthlyExpenseRoutes.DELETE("/:id", DeleteMonthlyExpense)
	monthlyExpenseRoutes.PUT("/", UpdateMonthlyExpense)
	monthlyExpenseRoutes.GET("/:id")
	monthlyExpenseRoutes.GET("/")
}

// ==================================================================
// 						CreateMonthlyExpense
// ==================================================================

type CreateMonthlyExpenseInput struct {
	Category string `json:"category" binding:"required"`
}

func CreateMonthlyExpense(context *gin.Context) {
	var input CreateMonthlyExpenseInput
	err := context.ShouldBindJSON(&input)
	if CheckError(context, err) {
		return
	}
	category := input.Category

	userId, err := utils.GetCurrentUserId(context)
	if CheckError(context, err) {
		return
	}

	monthlyExpense, err := database.CreateMonthlyExpense(userId, category)
	if CheckError(context, err) {
		return
	}

	context.JSON(http.StatusCreated, gin.H{"data": monthlyExpense})
}

// ==================================================================
// 						DeleteMonthlyExpense
// ==================================================================

func DeleteMonthlyExpense(context *gin.Context) {
	id, err := FetchId("id", context)
	if CheckError(context, err) {
		return
	}

	userId, err := utils.GetCurrentUserId(context)
	if CheckError(context, err) {
		return
	}

	monthlyExpense, err := database.GetMonthlyExpense(userId, uint(id))
	if CheckError(context, err) {
		return
	}

	err = monthlyExpense.Delete()
	if CheckError(context, err) {
		return
	}

	context.JSON(http.StatusOK, nil)
}

// ==================================================================
// 						UpdateMonthlyExpense
// ==================================================================

func UpdateMonthlyExpense(context *gin.Context) {
	var input database.MonthlyExpense
	err := context.ShouldBindJSON(&input)
	if CheckError(context, err) {
		return
	}

	userId, err := utils.GetCurrentUserId(context)
	if CheckError(context, err) {
		return
	}

	if input.UserID != userId {
		err = errors.New("this monthly expense does not belong to this user")
		CheckError(context, err)
		return
	}

	err = input.Update()
	if CheckError(context, err) {
		return
	}

	context.JSON(http.StatusOK, nil)
}

// ==================================================================
// 						GetMonthlyExpense
// ==================================================================

func GetMonthlyExpense(context *gin.Context) {
	id, err := FetchId("id", context)
	if CheckError(context, err) {
		return
	}

	userId, err := utils.GetCurrentUserId(context)
	if CheckError(context, err) {
		return
	}

	monthlyExpense, err := database.GetMonthlyExpense(userId, uint(id))
	if CheckError(context, err) {
		return
	}

	context.JSON(http.StatusCreated, gin.H{"data": monthlyExpense})
}

// ==================================================================
// 						GetAllMonthlyExpense
// ==================================================================

func GetAllMonthlyExpense(context *gin.Context) {
	userId, err := utils.GetCurrentUserId(context)
	if CheckError(context, err) {
		return
	}

	monthlyExpenses, err := database.GetAllMonthlyExpense(userId)
	if CheckError(context, err) {
		return
	}

	context.JSON(http.StatusCreated, gin.H{"data": monthlyExpenses})
}
