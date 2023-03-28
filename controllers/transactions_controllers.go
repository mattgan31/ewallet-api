package controllers

import (
	"ewallet-api/database"
	"ewallet-api/helpers"
	"ewallet-api/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

type TopupInput struct {
	Amount int `json:"amount" validate:"required"`
}

func Topup(c *gin.Context) {
	var Input TopupInput
	db := database.GetDB()
	User := models.User{}
	userData := c.MustGet("userData").(jwt.MapClaims)
	userID := uint(userData["id"].(float64))

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

	err = db.Debug().Model(&User).Where("id=?", userID).Update("balance", newBalance).Error

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	Transaction := models.Transaction{}
	Transaction.Amount = Input.Amount
	Transaction.Created_At = time.Now()
	Transaction.Updated_At = time.Now()
	Transaction.Description = ""
	Transaction.Type = "topup"
	Transaction.User_ID = int(userID)

	err = db.Debug().Create(&Transaction).Error
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
