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

type TransferInput struct {
	Amount           int    `json:"amount" validate:"required"`
	EmailDestination string `json:"email_destination" validate:"required"`
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
	description := "topup"
	HistoryTransaction(c, Input.Amount, userID.(uint), typeTransaction, description)
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
	description := "payment"
	HistoryTransaction(c, Input.Amount, userID.(uint), typeTransaction, description)
}

func Transfer(c *gin.Context) {
	db := database.GetDB()
	var Input TransferInput
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

	// newBalance2:=
	User2 := models.User{}

	err2 := db.Where("email=?", Input.EmailDestination).First(&User2).Error

	fmt.Println(User2)
	if err2 != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	fmt.Println(User2.Status)

	if User2.Status != "active" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Account Destination is inactive",
		})
		return
	}

	newBalance2 := User2.Balance + Input.Amount
	err2 = db.Model(&User2).Where("email=?", Input.EmailDestination).Update("balance", newBalance2).Error
	if err2 != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}
	newBalance := User.Balance - Input.Amount
	typeTransaction := "transfer"
	err = db.Debug().Model(&User).Where("id=?", userID).Update("balance", newBalance).Error

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}
	description := "transfer to " + Input.EmailDestination
	HistoryTransaction(c, Input.Amount, userID.(uint), typeTransaction, description)
}

func HistoryTransaction(c *gin.Context, amount int, userID uint, typeTransaction string, description string) {
	db := database.GetDB()
	Transaction := models.Transaction{}
	Transaction.Amount = amount
	Transaction.Created_At = time.Now()
	Transaction.Updated_At = time.Now()
	Transaction.Description = description
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
