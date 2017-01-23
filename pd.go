package main

import (
  "fmt"
  "io/ioutil"
  "log"
  "net/http"
  "encoding/json"
  "os"
)

type User struct {
    Id string `json:"id"`
    Summary string `json:"summary"`
    HtmlUrl string `json:"html_url"`
}

type Schedule struct {
    Id string `json:"id"`
    Summary string `json:"summary"`
    Name string `json:"name"`
    Users []User `json:"users"`
}

type PdResponse struct {
    Schedules []Schedule `json:"schedules"`
}

type SlackChannel struct {
    Id string `json:"id"`
    Name string `json:"name"`
    Topic SlackTopic `json:"topic"`
}

type SlackTopic struct {
    Value string `json:"value"`
    UnixTimeSet int64 `json:"last_set"`
    Creator string `json:"creator"`
}

type SlackResponse struct {
    Channel SlackChannel `json:"channel"`
    Ok bool `json:"ok"`
}


func getChannelTopic(slackToken string, channelId string) string {
  var info_url  = "https://slack.com/api/channels.info?token="+slackToken+"&channel="+channelId
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

func updateChannelTopic(slackToken string, topic string, channelId string) {

    //var topic_url = "https://slack.com/api/channels.setTopic?token="+slackToken+"&channel="+channelId+"&topic="+topic

    // TODO error if not present
}

func main() {
  var SLACK_TOKEN = os.Getenv("SLACK_TOKEN")

  request, _ := http.NewRequest("GET", "https://api.pagerduty.com/schedules", nil)
  request.Header.Set("Accept", "application/vnd.pagerduty+json;version=2")
  request.Header.Set("Authorization", "Token token=YqTrgWo4PMh5k5wZubyu")

  resp, err := http.DefaultClient.Do(request)
  if err != nil {
    log.Fatal(err)
  }

  body, _ := ioutil.ReadAll(resp.Body)

  var dat PdResponse

  if err := json.Unmarshal(body, &dat); err != nil {
       panic(err)
  }

  for _, schedule:= range dat.Schedules {
      fmt.Println(schedule)
      fmt.Println("")
  }

  fmt.Println(getChannelTopic(SLACK_TOKEN, "C11L5HUJY"))

//  schedules := dat["schedules"]

//  var oncalls map[string]string
//  oncalls = make(map[string]string)

//  fmt.Println(schedules)

//  fmt.Println(string(body))
}
