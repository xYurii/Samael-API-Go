package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"time"
)

type DiscordUser struct {
	ID            string `json:"id"`
	Username      string `json:"username"`
	Discriminator string `json:"discriminator"`
	Avatar        string `json:"avatar"`
	Verified      bool   `json:"verified"`
	Email         string `json:"email"`
	Flags         int    `json:"flags"`
	Banner        string `json:"banner"`
	AccentColor   int    `json:"accent_color"`
	PremiumType   int    `json:"premium_type"`
	PublicFlags   int    `json:"public_flags"`
}

type IPDetails struct {
	Status        string  `json:"status"`
	Continent     string  `json:"continent"`
	ContinentCode string  `json:"continentCode"`
	Country       string  `json:"country"`
	CountryCode   string  `json:"countryCode"`
	Region        string  `json:"region"`
	RegionName    string  `json:"regionName"`
	City          string  `json:"city"`
	District      string  `json:"district"`
	Zip           string  `json:"zip"`
	Lat           float64 `json:"lat"`
	Lon           float64 `json:"lon"`
	Timezone      string  `json:"timezone"`
	Offset        int     `json:"offset"`
	Currency      string  `json:"currency"`
	ISP           string  `json:"isp"`
	Org           string  `json:"org"`
	AS            string  `json:"as"`
	ASName        string  `json:"asname"`
	Reverse       string  `json:"reverse"`
	Mobile        bool    `json:"mobile"`
	Proxy         bool    `json:"proxy"`
	Hosting       bool    `json:"hosting"`
	Query         string  `json:"query"`
}

func RandomInt(min, max uint64) uint64 {
	if min >= max {
		return 0
	}
	rand.Seed(time.Now().UnixNano())
	return uint64(rand.Intn(int(max-min+1))) + min
}

func GetUser(TOKEN string) (DiscordUser, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://discord.com/api/v9/users/@me", nil)
	if err != nil {
		return DiscordUser{}, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", TOKEN))

	resp, err := client.Do(req)
	if err != nil {
		return DiscordUser{}, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return DiscordUser{}, err
	}

	var data DiscordUser
	err = json.Unmarshal(body, &data)
	if err != nil {
		return DiscordUser{}, err
	}

	return data, nil
}

func GetIPInfo(ip string) (IPDetails, error) {
	url := fmt.Sprintf("http://ip-api.com/json/%s?fields=66846719", ip)
	resp, err := http.Get(url)
	if err != nil {
		return IPDetails{}, err
	}
	defer resp.Body.Close()

	var ipInfo IPDetails
	err = json.NewDecoder(resp.Body).Decode(&ipInfo)
	if err != nil {
		return IPDetails{}, err
	}

	return ipInfo, nil
}

func ParseDuration(hour int) int64 {
	return int64((time.Duration(hour) * time.Hour) / time.Millisecond)
}

func InCommandCooldown(lastUsageTime int64, cooldownHours int) bool {
	cooldown := time.Duration(cooldownHours) * time.Hour / time.Millisecond
	currentTimeMilliseconds := time.Now().UnixNano() / int64(time.Millisecond)

	return currentTimeMilliseconds-lastUsageTime <= int64(cooldown)
}

func FormatTime(duration time.Duration, maxUnits int) string {
	days := int(duration.Hours() / 24)
	duration -= time.Duration(days) * 24 * time.Hour
	hours := int(duration.Hours())
	duration -= time.Duration(hours) * time.Hour
	minutes := int(duration.Minutes())
	duration -= time.Duration(minutes) * time.Minute
	seconds := int(duration.Seconds())

	units := []struct {
		value int
		label string
	}{
		{days, "dia"},
		{hours, "hora"},
		{minutes, "minuto"},
		{seconds, "segundo"},
	}

	result := ""
	count := 0
	for _, unit := range units {
		if unit.value > 0 && count < maxUnits {
			if result != "" {
				if count > 0 {
					if count == maxUnits-1 {
						result += " e "
					} else {
						result += ", "
					}
				}
			}
			result += fmt.Sprintf("%d %s", unit.value, unit.label)
			if unit.value != 1 {
				result += "s" // plural
			}
			count++
		}
	}

	return result
}
