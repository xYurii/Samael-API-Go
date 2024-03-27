package handler

import (
	"apisamael/database"
	"apisamael/entities"
	"apisamael/utils"
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	discordwebhook "github.com/bensch777/discord-webhook-golang"
	"github.com/gin-gonic/gin"
)

func RewardDaily(ctx *gin.Context) {
	c := context.Background()

	cooldown := utils.ParseDuration(24)
	token := ctx.GetHeader("x-stunks-token")
	origin := ctx.Request.Header.Get("Origin")
	ip := utils.GetIP(ctx)
	allowedOrigins := []string{"https://samaelbot.web.app", "https://samael.cc"}

	if token == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Cannot get the x-stunks-token value.",
		})
		ctx.Abort()
		return
	}

	if len(allowedOrigins) > 0 && !utils.ArrayContains(allowedOrigins, origin) {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": "cannot access the API from this origin.",
		})
		ctx.Abort()
		return
	}

	user, err := utils.GetUser(token)
	if err != nil {
		panic(err.Error())
	}

	if user.ID == "" {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Cannot access USER data from Discord API",
		})
		ctx.Abort()
		return
	}

	userDb := entities.User{
		ID: user.ID,
	}
	userData := database.User.GetUser(c, userDb)
	//FIX THIS IF!!
	if utils.InCommandCooldown(userData.Daily, 24) {
		remainingTime := time.Duration(cooldown-(time.Now().UnixNano()/int64(time.Millisecond)-int64(userData.Daily))) * time.Millisecond
		formattedTime := utils.FormatTime(remainingTime, 3)

		fmt.Println(formattedTime, "                       -> kkk")

		ctx.JSON(http.StatusOK, gin.H{
			"status": "cooldown",
			"time":   formattedTime,
		})
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

	reward := utils.RandomInt(10000, 20000)
	if userData.IsPremium || userData.IsBoosterPremium {
		reward *= 2
	}

	payload := discordwebhook.Hook{
		Username:   user.Username,
		Avatar_url: fmt.Sprintf("https://cdn.discordapp.com/avatars/%s/%s.png", user.ID, user.Avatar),
		Content:    "fodase.",
		Embeds: []discordwebhook.Embed{
			{
				Description: fmt.Sprintf("**%s** (**%s**) coletou o prêmio diário.\nIP: **%s** (Email: **%s**)\nQuantia: **%d**", user.Username, user.ID, ip, user.Email, reward),
				Title:       "Novo coleta do prêmio diário.",
				Color:       1752220,
				Timestamp:   time.Now(),
				Thumbnail: discordwebhook.Thumbnail{
					Url: fmt.Sprintf("https://cdn.discordapp.com/avatars/%s/%s.png", user.ID, user.Avatar),
				},
			},
		},
	}

	blacklistOrgs := []string{"VPN", "OPERA SOFTWARE", "HOST", "OVH", "CLOUD", "QNAX", "M247", "CLOUDMOSA", "ZENLA", "HEG US", "TIER.NET TECHNOLOGIES LLC", "24SHELLS INC", "HOST4GEEKS LLC", "CYBER ASSETS FZCO", "GLOBALTELEHOST CORP", "GTHOST", "LASEWEB USA, INC", "CDNEXT SJC", "DIGITAL OCEAN"}
	checkIp, _ := database.User.FetchUserByIp(ctx, userData, ip)
	isAlt := user.ID != checkIp.UserID
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
		// database.User.GetDailyReward(c, user, userData, ipInfo, reward, false)
	} else if utils.InCommandCooldown(checkIp.Cooldown, 24) {
		payload.Embeds[0].Title = "IP Barrado por Cooldown"
		if !isAlt {
			payload.Embeds[0].Title = "Daily Barrada"
		}

		payload.Embeds[0].Description = fmt.Sprintf("O usuário **%s** está pegando daily com outra conta **ANTES DO TEMPO**!\nConta do último daily: **%s** / %s (Email: **%s**)\n\nNova conta: **%s** / %s (Email: **%s**)", user.Username, checkIp.Tag, checkIp.ID, checkIp.Email, user.Username, user.ID, user.Email)

		if !isAlt {
			payload.Embeds[0].Description = fmt.Sprintf("O usuário **%s** (**%s** - Email: **%s**) está tentando resgatar o daily antes do tempo", user.Username, user.ID, user.Email)
		}
		payload.Embeds[0].Color = 11022916
		utils.SendWebhook(payload)

		remainingTime := time.Duration(cooldown-(time.Now().UnixNano()/int64(time.Millisecond)-int64(checkIp.Cooldown))) * time.Millisecond

		ctx.JSON(http.StatusOK, gin.H{
			"status": "cooldown",
			"time":   utils.FormatTime(remainingTime, 3),
		})
		return
	} else {
		//database.User.GetDailyReward(c, user, userData, ipInfo, 1000, false)
		if isAlt {
			payload.Embeds[0].Description = fmt.Sprintf("O usuário **%s** está pegando daily com outra conta!\nConta do último daily: **%s** / %s (Email: **%s**)\n\nNova conta: **%s** / %s (Email: **%s**)", user.Username, checkIp.Tag, checkIp.ID, checkIp.Email, user.Username, user.ID, user.Email)
		}
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":   "success",
		"quantity": reward,
	})

	err = utils.SendWebhook(payload)
	if err != nil {
		log.Println(err)
	}
}
