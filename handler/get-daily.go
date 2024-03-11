package handler

import (
	"apisamael/database"
	"apisamael/entities"
	"apisamael/utils"
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func RewardDaily(ctx *gin.Context) {
	c := context.Background()

	token := ctx.Query("x-stunks-token")
	ip := getIp(ctx)
	allowedOrigins := []string{}
	origin := ctx.GetHeader("Origin")

	if token == "" {
		return
	}

	if len(allowedOrigins) > 0 && contains(allowedOrigins, origin) {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status": "error",
		})
		ctx.Abort()
		return
	}

	user, err := GetUser(token)
	if err != nil {
		panic(err.Error())
	}
	userDb := entities.User{
		ID: user.ID,
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

	isBlacklisted, _ := database.User.IsBlacklisted(c, userData)
	if isBlacklisted {
		ctx.JSON(http.StatusForbidden, gin.H{
			"status": "blacklist",
		})
		ctx.Abort()
		return
	}

	blacklistOrgs := []string{"VPN", "OPERA SOFTWARE", "HOST", "OVH", "CLOUD", "QNAX", "M247", "CLOUDMOSA", "ZENLA", "HEG US", "TIER.NET TECHNOLOGIES LLC", "24SHELLS INC", "HOST4GEEKS LLC", "CYBER ASSETS FZCO", "GLOBALTELEHOST CORP", "GTHOST", "LASEWEB USA, INC", "CDNEXT SJC", "DIGITAL OCEAN"}
	fmt.Println(blacklistOrgs)
	checkIp, _ := database.User.FetchUserByIp(ctx, userData, ip)
	isAlt := user.ID != checkIp.UserID
	fmt.Println(isAlt)
	ipInfo, _ := utils.GetIPInfo(ip)
	userData.UserTasks.Daily = true

	if checkIp.ID == "" {
		for _, org := range blacklistOrgs {
			if strings.Contains(strings.ToUpper(ipInfo.Org), strings.ToUpper(org)) {
				ctx.JSON(http.StatusForbidden, gin.H{
					"status": "vpn",
				})
				ctx.Abort()
				return
			}
		}
		database.User.GetDailyReward(c, user, userData, ipInfo, 1000, false)
	}

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
