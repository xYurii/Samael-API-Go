package handler

import (
	"apisamael/database"
	"apisamael/entities"
	"apisamael/utils"
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func RewardDaily(ctx *gin.Context) {
	c := context.Background()

	token := ctx.Query("x-stunks-token")
	if token == "" {
		return
	}

	allowedOrigins := []string{}
	origin := ctx.GetHeader("Origin")
	if len(allowedOrigins) > 0 && contains(allowedOrigins, origin) {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status": "error",
		})
		ctx.Abort()
		return
	}

	ip := getIp(ctx)
	fmt.Println(ip) //just for lint ignore

	user, err := GetUser(token)
	if err != nil {
		panic(err.Error())
	}
	userDb := entities.User{
		ID: user["id"].(string),
	}
	userData := database.User.GetUser(c, userDb)
	if utils.InCommandCooldown(int64(userData.Daily), 24) {
		cooldown := utils.ParseDuration(24)
		remainingTime := time.Duration(cooldown-(time.Now().UnixNano()/int64(time.Millisecond)-int64(userData.Daily))) * time.Millisecond

		ctx.JSON(http.StatusOK, gin.H{
			"status": "cooldown",
			"time":   utils.FormatTime(remainingTime, 3),
		})
		ctx.Abort()
		return
	}

	isBlacklisted := database.User.IsBlacklisted(c, userData)
	if isBlacklisted {
		ctx.JSON(http.StatusForbidden, gin.H{
			"status": "blacklist",
		})
		ctx.Abort()
		return
	}

	//continue tomorrow

}

func getIp(ctx *gin.Context) string {
	ip := ctx.GetHeader("x-forwarded-for")
	if ip == "" {
		ip = ctx.ClientIP()
	}
	return ip
}

func contains(arr []string, str string) bool {
	for _, a := range arr {
		if a == str {
			return true
		}
	}
	return false
}
