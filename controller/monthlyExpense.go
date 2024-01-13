package controller

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/SgtMilk/fin-planning-backend/database"
	"github.com/SgtMilk/fin-planning-backend/utils"
	"github.com/gin-gonic/gin"
)

func CreateMonthlyExpenseRouter(router *gin.RouterGroup) {
	monthlyExpenseRoutes := router.Group("/monthly_expense")
	monthlyExpenseRoutes.POST("/save", SaveMonthlyExpense)
	monthlyExpenseRoutes.POST("/", CreateMonthlyExpense)
	monthlyExpenseRoutes.DELETE("/:id", DeleteMonthlyExpense)
	monthlyExpenseRoutes.PUT("/", UpdateMonthlyExpense)
	monthlyExpenseRoutes.GET("/:id")
	monthlyExpenseRoutes.GET("/")
}

// ==================================================================
// 						CreateMonthlyExpense
// ==================================================================

func SaveMonthlyExpense(context *gin.Context) {
	var inputs []database.MonthlyExpense
	err := context.ShouldBindJSON(&inputs)
	if CheckError(context, err) {
		return
	}

	userId, err := utils.GetCurrentUserId(context)
	if CheckError(context, err) {
		return
	}
	user, err := database.GetUserWithMontlyExpenses(userId)
	if CheckError(context, err) {
		return
	}

	// classifying incoming monthlyExpenses
	var newMonthlyExpenses, oldMonthlyExpenses, deletedMonthlyExpenses []database.MonthlyExpense
	for _, monthlyExpense := range inputs {
		if monthlyExpense.UserID != user.ID {
			err := errors.New("monthly expense ID#" + strconv.FormatUint(uint64(monthlyExpense.ID), 10) + " does not belong to this user")
			CheckError(context, err)
			return
		}

		if monthlyExpense.ID == 0 {
			newMonthlyExpenses = append(newMonthlyExpenses, monthlyExpense)
		} else {
			oldMonthlyExpenses = append(oldMonthlyExpenses, monthlyExpense)
		}
	}

	for _, oldMonthlyExpense := range user.MonthlyExpenses {
		if !arrContains(inputs, oldMonthlyExpense.ID) {
			deletedMonthlyExpenses = append(deletedMonthlyExpenses, oldMonthlyExpense)
		}
	}

	err = database.DeleteMonthlyExpenses(deletedMonthlyExpenses)
	if CheckError(context, err) {
		return
	}

	err = database.CreateMonthlyExpenses(newMonthlyExpenses)
	if CheckError(context, err) {
		return
	}

	err = database.SaveMonthlyExpenses(oldMonthlyExpenses)
	if CheckError(context, err) {
		return
	}

	err = user.UpdateMonthlyExpenses(append(oldMonthlyExpenses, newMonthlyExpenses...))
	if CheckError(context, err) {
		return
	}

	context.JSON(http.StatusOK, nil)
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

// ==================================================================
// 								HELPERS
// ==================================================================

func arrContains(arr []database.MonthlyExpense, value uint) bool {
	for _, val := range arr {
		if val.ID == value {
			return true
		}
	}
	return false
}
