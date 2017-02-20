package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jasonlvhit/gocron"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
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
	var infoUrl = "https://slack.com/api/channels.info"

	request, _ := http.NewRequest("GET", infoUrl, nil)
	var query = request.URL.Query()
	query.Add("token", slackToken)
	query.Add("channel", channelId)
	request.URL.RawQuery = query.Encode()

	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		log.Fatal(err)
	}

	body, _ := ioutil.ReadAll(resp.Body)

	var dat SlackResponse
	if err := json.Unmarshal(body, &dat); err != nil {
		panic(err)
	}

	fmt.Println("Current topic for channel " + channelId + ": " + dat.Channel.Topic.Value)
	return dat.Channel.Topic.Value
}

func updateChannelTopic(slackToken string, topic string, channelId string) {

	var updateUrl = "https://slack.com/api/channels.setTopic"

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

func getOncallName(pdScheduleId string) (string, error) {
	var PD_TOKEN = os.Getenv("PD_TOKEN")

	// https://api.pagerduty.com/oncalls?time_zone=UTC&schedule_ids%5B%5D=P7CMRA9%2CTESTETS
	request, _ := http.NewRequest("GET", "https://api.pagerduty.com/oncalls", nil)
	var query = request.URL.Query()
	query.Add("schedule_ids[]", pdScheduleId)
	request.URL.RawQuery = query.Encode()
	request.Header.Set("Accept", "application/vnd.pagerduty+json;version=2")
	request.Header.Set("Authorization", "Token token="+PD_TOKEN)

	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		return "", err
	}

	body, _ := ioutil.ReadAll(resp.Body)

	var dat PdOncallsResponse

	if err := json.Unmarshal(body, &dat); err != nil {
		return "", err
	}

	for _, oncall := range dat.Oncalls {
		if oncall.Schedule.Id == pdScheduleId {
			return oncall.User.Summary, nil
		}
	}

	return "", errors.New("couldn't get the on call name")
}

func getOncallAndUpdateSlackChannel(slackChannelId string, pdScheduleId string) {
	fmt.Println("Checking for slack channel with id " + slackChannelId)
	var slackToken = os.Getenv("SLACK_TOKEN")
	var prefix = "Engineer on call: "
	var currentTopic = getChannelTopic(slackToken, slackChannelId)
	var oncallName, err = getOncallName(pdScheduleId)
	if err != nil {
		fmt.Println(err)
		return
	}

	var topic = strings.TrimSpace(prefix + oncallName)
	if currentTopic != "" && currentTopic != topic {
		fmt.Println("Setting the new topic: " + topic)
		updateChannelTopic(slackToken, prefix+oncallName, slackChannelId)
	}
}

func main() {
	var PLATFORM_SCHEDULE_ID = "P7CMRA9"
	var PLATFORM_CHANNEL_ID = "C11L5HUJY"
	fmt.Println("Starting the bot...")

	// "G2K8LQ3SA"
	gocron.Every(1).Day().At("19:05").Do(getOncallAndUpdateSlackChannel, PLATFORM_CHANNEL_ID, PLATFORM_SCHEDULE_ID)

	//<-gocron.Start()
}
