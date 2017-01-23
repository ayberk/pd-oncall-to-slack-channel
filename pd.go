package main

import (
	"encoding/json"
	"fmt"
	"github.com/jasonlvhit/gocron"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type PdUser struct {
	Id      string `json:"id"`
	Summary string `json:"summary"`
	HtmlUrl string `json:"html_url"`
}

type PdSchedule struct {
	Id      string   `json:"id"`
	Summary string   `json:"summary"`
	Name    string   `json:"name"`
	Users   []PdUser `json:"users"`
}

type PdOncall struct {
	Schedule PdSchedule `json:"schedule"`
	User     PdUser     `json:"user"`
}

type PdSchedulesResponse struct {
	Schedules []PdSchedule `json:"schedules"`
}

type PdOncallsResponse struct {
	Oncalls []PdOncall `json:"oncalls"`
}

type SlackChannel struct {
	Id    string     `json:"id"`
	Name  string     `json:"name"`
	Topic SlackTopic `json:"topic"`
}

type SlackTopic struct {
	Value       string `json:"value"`
	UnixTimeSet int64  `json:"last_set"`
	Creator     string `json:"creator"`
}

type SlackResponse struct {
	Channel SlackChannel `json:"channel"`
	Ok      bool         `json:"ok"`
}

func getChannelTopic(slackToken string, channelId string) string {
	var info_url = "https://slack.com/api/channels.info?token=" + slackToken + "&channel=" + channelId
	request, _ := http.NewRequest("GET", info_url, nil)

	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		log.Fatal(err)
	}

	body, _ := ioutil.ReadAll(resp.Body)

	var dat SlackResponse
	if err := json.Unmarshal(body, &dat); err != nil {
		panic(err)
	}

	return dat.Channel.Topic.Value
}

// platform primary schedule id = "P7CMRA9"

func updateChannelTopic(slackToken string, topic string, channelId string) {

	var currentTopic = getChannelTopic(slackToken, channelId)
	var updateUrl = "https://slack.com/api/groups.setTopic"

	if topic != currentTopic {
		request, _ := http.NewRequest("GET", updateUrl, nil)
		var query = request.URL.Query()
		query.Add("token", slackToken)
		query.Add("channel", channelId)
		query.Add("topic", topic)
		request.URL.RawQuery = query.Encode()

		_, err := http.DefaultClient.Do(request)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func getOncallName(pdScheduleId string) string {
	var PD_TOKEN = os.Getenv("PD_TOKEN")

	request, _ := http.NewRequest("GET", "https://api.pagerduty.com/oncalls", nil)
	request.Header.Set("Accept", "application/vnd.pagerduty+json;version=2")
	request.Header.Set("Authorization", "Token token="+PD_TOKEN)

	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		log.Fatal(err)
	}

	body, _ := ioutil.ReadAll(resp.Body)

	var dat PdOncallsResponse

	if err := json.Unmarshal(body, &dat); err != nil {
		panic(err)
	}

	for _, oncall := range dat.Oncalls {
		if oncall.Schedule.Id == pdScheduleId {
			fmt.Println(oncall.User.Summary)
			return oncall.User.Summary
		}
	}

	// TODO return err
	return ""
}

func getOncallAndUpdateSlackChannel(slackChannelId string, pdScheduleId string) {
	var SLACK_TOKEN = os.Getenv("SLACK_TOKEN")

	var oncallName = getOncallName(pdScheduleId)

	//"C11L5HUJY" -> platformChannelId
	// Don't forget to change API endpoint from groups to channels!
	updateChannelTopic(SLACK_TOKEN, "On call: "+oncallName, slackChannelId)
}
func main() {
	var PLATFORM_SCHEDULE_ID = "P7CMRA9"

	gocron.Every(1).Day().At("22:28").Do(getOncallAndUpdateSlackChannel, "G2K8LQ3SA", PLATFORM_SCHEDULE_ID)
	<-gocron.Start()
}
