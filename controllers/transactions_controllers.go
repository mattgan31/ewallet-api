package controllers

import (
	"ewallet-api/database"
	"ewallet-api/helpers"
	"ewallet-api/models"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type TransactionInput struct {
	Amount int `json:"amount" validate:"required"`
}

func Topup(c *gin.Context) {
	var Input TransactionInput
	db := database.GetDB()
	User := models.User{}

	userID, exists := c.Get("user_id")

	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "User Not Found",
		})
	}

	contentType := helpers.GetContentType(c)
	if contentType == appJSON {
		c.ShouldBindJSON(&Input)
	} else {
		c.ShouldBind(&Input)
	}

	err := db.First(&User, userID).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	newBalance := User.Balance + Input.Amount
	typeTransaction := "topup"
	err = db.Debug().Model(&User).Update("balance", newBalance).Error

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}
	HistoryTransaction(c, Input.Amount, userID.(uint), typeTransaction)
}

func Payment(c *gin.Context) {
	var Input TransactionInput
	db := database.GetDB()
	User := models.User{}
	userID, exists := c.Get("user_id")

	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "User Not Found",
		})
	}

	contentType := helpers.GetContentType(c)
	if contentType == appJSON {
		c.ShouldBindJSON(&Input)
	} else {
		c.ShouldBind(&Input)
	}

	err := db.First(&User, userID).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	if User.Balance < Input.Amount {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Your Balance is not enough for this payment",
		})
		return
	}

	newBalance := User.Balance - Input.Amount
	typeTransaction := "payment"
	err = db.Debug().Model(&User).Where("id=?", userID).Update("balance", newBalance).Error

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}
	HistoryTransaction(c, Input.Amount, userID.(uint), typeTransaction)
}

func HistoryTransaction(c *gin.Context, amount int, userID uint, typeTransaction string) {
	db := database.GetDB()
	Transaction := models.Transaction{}
	Transaction.Amount = amount
	Transaction.Created_At = time.Now()
	Transaction.Updated_At = time.Now()
	Transaction.Description = ""
	Transaction.Type = typeTransaction
	Transaction.User_ID = int(userID)

	err := db.Debug().Create(&Transaction).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   &Transaction,
	})
}

func GetHistory(c *gin.Context) {
	db := database.GetDB()
	Transaction := []models.Transaction{}
	var result gin.H
	userID, exists := c.Get("user_id")

	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "User Not Found",
		})
		return
	}

	fmt.Println(userID)

	err := db.Where("user_id=?", userID).Find(&Transaction).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Data Not Found",
		})
		return
	}

	if len(Transaction) <= 0 {
		result = gin.H{
			"result": "data not found",
		}
	} else {
		result = gin.H{
			"result": Transaction,
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   result,
	})
}
