package handler

import (
	"apisamael/utils"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func GetUser(TOKEN string) (utils.DiscordUser, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://discord.com/api/v9/users/@me", nil)
	if err != nil {
		return utils.DiscordUser{}, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", TOKEN))

	resp, err := client.Do(req)
	if err != nil {
		return utils.DiscordUser{}, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return utils.DiscordUser{}, err
	}

	var data utils.DiscordUser
	err = json.Unmarshal(body, &data)
	if err != nil {
		return utils.DiscordUser{}, err
	}

	return data, nil
}
