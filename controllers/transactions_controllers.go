package controllers

import (
	"ewallet-api/database"
	"ewallet-api/helpers"
	"ewallet-api/models"
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
	var input TransactionInput
	db := database.GetDB()
	var user models.User

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "User Not Found",
		})
		return
	}

	contentType := helpers.GetContentType(c)
	if contentType == appJSON {
		c.ShouldBindJSON(&input)
	} else {
		c.ShouldBind(&input)
	}

	if err := db.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	newBalance := user.Balance + input.Amount
	typeTransaction := "topup"

	if err := db.Debug().Model(&user).Update("balance", newBalance).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	description := "topup"
	HistoryTransaction(c, input.Amount, userID.(uint), typeTransaction, description)

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"user_id": userID,
			"balance": newBalance,
		},
	})
}

func Payment(c *gin.Context) {
	var Input TransactionInput
	db := database.GetDB()
	userID, exists := c.Get("user_id")

	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "User Not Found",
		})
		return
	}

	contentType := helpers.GetContentType(c)
	if contentType == appJSON {
		c.ShouldBindJSON(&Input)
	} else {
		c.ShouldBind(&Input)
	}

	User := models.User{GormModel: models.GormModel{ID: userID.(uint)}}

	if err := db.First(&User).Error; err != nil {
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

	if err := db.Model(&User).Update("balance", newBalance).Error; err != nil {
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

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "User Not Found",
		})
		return
	}

	var Input TransferInput
	contentType := helpers.GetContentType(c)
	if contentType == appJSON {
		c.ShouldBindJSON(&Input)
	} else {
		c.ShouldBind(&Input)
	}

	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var User, User2 models.User
	err := tx.Where("id = ?", userID).First(&User).Error
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	if User.Balance < Input.Amount {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Your Balance is not enough for this payment",
		})
		return
	}

	err = tx.Where("email = ?", Input.EmailDestination).First(&User2).Error
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	if User2.Status != "active" {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Account Destination is inactive",
		})
		return
	}

	newBalance2 := User2.Balance + Input.Amount
	err = tx.Model(&User2).Where("email = ?", Input.EmailDestination).Update("balance", newBalance2).Error
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	newBalance := User.Balance - Input.Amount
	err = tx.Model(&User).Where("id = ?", userID).Update("balance", newBalance).Error
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	typeTransaction := "transfer"
	description := "transfer to " + Input.EmailDestination
	HistoryTransaction(c, Input.Amount, userID.(uint), typeTransaction, description)

	tx.Commit()

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Transfer Success",
	})
}

func HistoryTransaction(c *gin.Context, amount int, userID uint, typeTransaction, description string) {
	db := database.GetDB()

	transaction := models.Transaction{
		Amount: amount,
		GormModel: models.GormModel{
			Created_At: time.Now(),
			Updated_At: time.Now(),
		},
		Description: description,
		Type:        typeTransaction,
		User_ID:     int(userID),
	}

	if err := db.Debug().Create(&transaction).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   &transaction,
	})
}

func GetHistory(c *gin.Context) {
	db := database.GetDB()

	// Mendapatkan ID pengguna dari context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "User Not Found",
		})
		return
	}

	// Mencari transaksi menggunakan ID pengguna
	var transactions []models.Transaction
	if err := db.Where("user_id = ?", userID).Find(&transactions).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Data Not Found",
		})
		return
	}

	// Menyiapkan data untuk dikembalikan dalam respons
	data := gin.H{
		"result": "data not found",
	}
	if len(transactions) > 0 {
		data = gin.H{
			"result": transactions,
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   data,
	})
}
