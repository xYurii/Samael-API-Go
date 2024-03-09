package router

import (
	"apisamael/handler"

	"github.com/gin-gonic/gin"
)

func initializeRoutes(router *gin.Engine) {
	basePath := "/api/v1"
	v1 := router.Group(basePath)
	{
		v1.GET("/rewardDaily", handler.RewardDaily)
	}
}
