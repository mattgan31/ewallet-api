package router

import (
	"ewallet-api/controllers"
	"ewallet-api/middleware"

	"github.com/gin-gonic/gin"
)

func StartServer() *gin.Engine {
	router := gin.Default()

	userRouter := router.Group("/user")
	{
		userRouter.POST("/register", controllers.UserRegister)
		userRouter.POST("/login", controllers.UserLogin)
		userRouter.GET("/detail", middleware.Authentication(), controllers.GetDetailUser)
	}
	transactionRouter := router.Group("/transaction")
	{
		transactionRouter.POST("/topup", middleware.Authentication(), controllers.Topup)
		transactionRouter.POST("/payment", middleware.Authentication(), controllers.Payment)
	}
	return router
}
