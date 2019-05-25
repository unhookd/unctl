package config

import (
	"fmt"
	"github.com/unhookd/unctl/lib"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

func NotifyCommunicationChannel(notificationConfiguration NotificationsTable) {
	if notificationConfiguration["provider"] == "slack" {
		slack_api_url := lib.GetEnv("SLACK_API_URL", "https://hooks.slack.com/services/")
		sn, message, responseBody := notifySlack(slack_api_url, notificationConfiguration["channel"], notificationConfiguration["text"])
		if !sn {
			log.Printf("Slack Notification Failed ## Message: %s, Response: %s", message, responseBody)
		} else {
			log.Printf("Slack Notification Succeeded ## Message: %s, Response: %s", message, responseBody)
		}
	}
}

func notifySlack(slack_api_url, channel, text string) (success bool, message, responseBody string) {
	message = fmt.Sprintf(`{"channel": "%v", "username": "unctl", "text": "%v", "icon_emoji": ":ops:"}`, channel, text)

	values := url.Values{}
	values.Add("payload", message)

	slack_secret := lib.GetEnv("SLACK_SECRET_KEY", "foobar-baz")
	endpoint := fmt.Sprintf("%s/%s", slack_api_url, slack_secret)
	req, err := http.NewRequest("POST", endpoint, strings.NewReader(values.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("api_url: %s, channel: %v slack_secret: %s", endpoint, channel, slack_secret)
		log.Printf("Unable to send POST request to Slack: %v", err.Error())
	} else {
		if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
			success = true
		} else {
			body, _ := ioutil.ReadAll(resp.Body)
			log.Printf("Slack responded with status %s and body: %s", resp.Status, string(body))
		}
	}

	defer resp.Body.Close()
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Unable to read response body %v", err.Error())
	}
	responseBody = string(bodyBytes)
	return success, message, responseBody
}
