package main

import (
  "fmt"
  "io/ioutil"
  "log"
  "net/http"
  "encoding/json"
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

func main() {
  request, _ := http.NewRequest("GET", "https://api.pagerduty.com/schedules", nil)
  request.Header.Set("Accept", "application/vnd.pagerduty+json;version=2")
  request.Header.Set("Authorization", "Token token=YqTrgWo4PMh5k5wZubyu")

  resp, err := http.DefaultClient.Do(request)
  if err != nil {
    log.Fatal(err)
  }

  body, _ := ioutil.ReadAll(resp.Body)

//  var dat map[string]interface{}

//  if err := json.Unmarshal(body, &dat); err != nil {
//       panic(err)
//   }

  var ss []Schedule

  if err := json.Unmarshal(body, &ss); err != nil {
       panic(err)
   }

  fmt.Println(ss)

//  schedules := dat["schedules"]

//  var oncalls map[string]string
//  oncalls = make(map[string]string)

//  fmt.Println(schedules)

//  fmt.Println(string(body))
}
