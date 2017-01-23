package main

import (
	"encoding/json"
	"fmt"
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
		fmt.Println(request.URL)

		resp, err := http.DefaultClient.Do(request)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(resp)
	}
}

func main() {
	var SLACK_TOKEN = os.Getenv("SLACK_TOKEN")
	var PD_TOKEN = os.Getenv("PD_TOKEN")
	var PLATFORM_SCHEDULE_ID = "P7CMRA9"

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

	var oncallName string
	for _, oncall := range dat.Oncalls {
		if oncall.Schedule.Id == PLATFORM_SCHEDULE_ID {
			fmt.Println(oncall.User.Summary)
			oncallName = oncall.User.Summary
			break
		}
	}

	//"C11L5HUJY" -> platformChannelId
	updateChannelTopic(SLACK_TOKEN, "On call: "+oncallName, "G2K8LQ3SA")
}
