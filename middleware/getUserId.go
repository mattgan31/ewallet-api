package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func GetUserID() gin.HandlerFunc {
	return func(c *gin.Context) {
		// dapatkan user ID dari request
		userData := c.MustGet("userData").(jwt.MapClaims)
		userID := uint(userData["id"].(float64))

		// jika user ID kosong, kembalikan error
		if int(userID) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Missing user ID"})
			return
		}

		// tambahkan user ID ke context
		c.Set("user_id", userID)

		// lanjutkan eksekusi middleware berikutnya atau fungsi utama
		c.Next()
	}
}
