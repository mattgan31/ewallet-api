package controllers

import (
	"ewallet-api/database"
	"ewallet-api/helpers"
	"ewallet-api/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/golang-jwt/jwt"
)

var (
	appJSON = "application/json"
)

func UserRegister(c *gin.Context) {
	db := database.GetDB()
	// contentType := helpers.GetContentType(c)

	user := models.User{}

	err := c.ShouldBindBodyWith(&user, binding.JSON)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	user.Balance = 0
	user.Status = "active"
	user.Created_At = time.Now()
	user.Updated_At = time.Now()

	err = db.Debug().Create(&user).Error
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
			"id":       user.ID,
			"fullname": user.Full_Name,
			"email":    user.Email,
		},
	})
}

func UserLogin(c *gin.Context) {
	db := database.GetDB()
	contentType := helpers.GetContentType(c)

	login := struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}{}

	if contentType == appJSON {
		c.ShouldBindJSON(&login)
	} else {
		c.ShouldBind(&login)
	}

	user := models.User{}
	// password := login.Password

	err := db.Debug().Where("email=?", login.Email).Take(&user).Error
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": "Invalid Email/Password",
		})
		return
	}

	comparePass := helpers.ComparePass([]byte(user.Password), []byte(login.Password))

	if !comparePass {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": "Invalid email/password",
		})
		return
	}

	token, err := helpers.GenerateToken(user.ID, user.Email)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

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

	user, err := models.GetUserByID(db, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"full_name": user.Full_Name,
			"email":     user.Email,
			"balance":   user.Balance,
			"status":    user.Status,
		},
	})
}
