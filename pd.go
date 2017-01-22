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

func updateChannelTopic(topic string, channelId string) {
    const SLACK_TOKEN = os.Getenv("SLACK_TOKEN")

    const var info_url  = "https://slack.com/api/channels.info?token="+SLACK_TOKEN+"&channel="+channelId
    const var topic_url = "https://slack.com/api/channels.setTopic?token="+SLACK_TOKEN+"&channel="+channelId+"&topic="+topic

    // TODO error if not present
}

func main() {
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

//  schedules := dat["schedules"]

//  var oncalls map[string]string
//  oncalls = make(map[string]string)

//  fmt.Println(schedules)

//  fmt.Println(string(body))
}
