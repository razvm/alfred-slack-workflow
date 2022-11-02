package main

import (
	"errors"
	aw "github.com/deanishe/awgo"
	"github.com/slack-go/slack"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

func removeEmptyStrings(s []string) []string {
	var r []string
	for _, str := range s {
		if str != "" {
			r = append(r, str)
		}
	}
	return r
}
func downloadFile(URL, fileName string) error {
	//Get the response bytes from the url
	response, err := http.Get(URL)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return errors.New("received non 200 response code")
	}
	//Create an empty file
	file, err := os.Create(cache_dir + "/" + fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	//Write the bytes to the file
	_, err = io.Copy(file, response.Body)
	if err != nil {
		return err
	}

	return nil
}
func updateChannels() {
	wf.NewItem("Update Channels").Valid(true)

	c := aw.NewCache(cache_dir)
	cfg := aw.NewConfig()
	token := cfg.Get("SLACK_TOKEN")
	api := slack.New(token)
	params := slack.GetConversationsParameters{}
	channels, _, err_channels := api.GetConversations(&params)
	team, err_team := api.GetTeamInfo()

	if err_channels != nil || err_team != nil {
		wf.Warn("Error", "Error occurred in Slack API ")
	}

	all_channels := make([]Channel, 0)
	for _, channel := range channels {
		all_channels = append(all_channels, Channel{
			Name:   channel.Name,
			ID:     channel.ID,
			TeamID: team.ID,
		})
	}

	users, err_users := api.GetUsers()
	if err_users != nil {
		wf.Warn("Error", "Error occurred in Slack API [users]")
	}

	for _, user := range users {
		if !user.Deleted && !user.IsBot {
			all_channels = append(all_channels, Channel{
				Name:   strings.Join(removeEmptyStrings([]string{user.RealName, user.Name, user.Profile.DisplayName}), " / "),
				ID:     user.ID,
				TeamID: team.ID,
			})
			err := downloadFile(user.Profile.Image192, user.ID+".png")
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	c.StoreJSON(cache_file, all_channels)
	wf.SendFeedback()
}
