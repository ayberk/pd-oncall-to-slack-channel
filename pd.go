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

	fmt.Println(string(body[:]))
	fmt.Println(dat.Channel.Topic.Value)

	return ""
}

// platform primary schedule id = "P7CMRA9"

func updateChannelTopic(slackToken string, topic string, channelId string) {

	var updateUrl = "https://slack.com/api/channels.setTopic?token=" + slackToken + "&channel=" + channelId + "&topic=" + topic

	request, _ := http.NewRequest("GET", updateUrl, nil)

	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(resp)
}

func main() {
	var SLACK_TOKEN = os.Getenv("SLACK_TOKEN")
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
		if oncall.Schedule.Id == "P7CMRA9" {
			fmt.Println(oncall.User.Summary)
			break
		}
	}

	fmt.Println(getChannelTopic(SLACK_TOKEN, "C11L5HUJY"))

	//  schedules := dat["schedules"]

	//  var oncalls map[string]string
	//  oncalls = make(map[string]string)

	//  fmt.Println(schedules)

	//  fmt.Println(string(body))
}
