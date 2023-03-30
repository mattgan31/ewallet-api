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
		transactionRouter.POST("/topup", middleware.Authentication(), middleware.GetUserID(), controllers.Topup)
		transactionRouter.POST("/payment", middleware.Authentication(), middleware.GetUserID(), controllers.Payment)
		transactionRouter.GET("/history", middleware.Authentication(), middleware.GetUserID(), controllers.GetHistory)
		transactionRouter.POST("/transfer", middleware.Authentication(), middleware.GetUserID(), controllers.Transfer)
	}
	return router
}
