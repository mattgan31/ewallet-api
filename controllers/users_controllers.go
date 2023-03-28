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

var (
	appJSON = "application/json"
)

type Login struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func UserRegister(c *gin.Context) {
	db := database.GetDB()
	contentType := helpers.GetContentType(c)

	User := models.User{}

	if contentType == appJSON {
		c.ShouldBindJSON(&User)
	} else {
		c.ShouldBind(&User)
	}

	User.Balance = 0
	User.Status = "active"
	User.Created_At = time.Now()
	User.Updated_At = time.Now()

	err := db.Debug().Create(&User).Error

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{
		"status": "success",
		"data": gin.H{
			"id":       User.ID,
			"fullname": User.Full_Name,
			"email":    User.Email,
		},
	})
}

func UserLogin(c *gin.Context) {
	var Login Login
	db := database.GetDB()
	contentType := helpers.GetContentType(c)

	User := models.User{}
	password := ""
	if contentType == appJSON {
		c.ShouldBindJSON(&Login)
	} else {
		c.ShouldBind(&Login)
	}

	password = Login.Password

	err := db.Debug().Where("email=?", Login.Email).Take(&User).Error
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": "Invalid Email/Password",
		})
		return
	}

	comparePass := helpers.ComparePass([]byte(User.Password), []byte(password))

	if !comparePass {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": "Invalid email/password",
		})
		return
	}

	token := helpers.GenerateToken(User.ID, User.Email)

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"token": token,
		},
	})
}

func GetDetailUser(c *gin.Context) {
	db := database.GetDB()
	userData := c.MustGet("userData").(jwt.MapClaims)
	userID := uint(userData["id"].(float64))
	User := models.User{}

	err := db.First(&User, userID).Error

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{"full_name": User.Full_Name,
			"email":   User.Email,
			"balance": User.Balance,
			"status":  User.Status},
	})
}
