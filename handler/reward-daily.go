package handler

import (
	"github.com/gin-gonic/gin"
)

func RewardDaily(ctx *gin.Context) {
	token := ctx.Query("token")
	if token == "" {
		return
	}
}
