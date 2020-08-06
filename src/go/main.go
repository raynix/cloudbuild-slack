package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"cloud.google.com/go/pubsub"
)

var projectID = "idyllic-depth-239301"
var topicName = "cloud-builds"
var subName = "cloud-build-slack-sub"

type slackMessage struct {
	Text string `json:"text"`
}

func hello() (helloOutput string) {
	return "Hello, world!"
}

func findOrCreateSub(ctx context.Context, client *pubsub.Client, topic *pubsub.Topic, subName string) (sub *pubsub.Subscription) {
	sub = client.Subscription(subName)
	ok, err := sub.Exists(ctx)
	if err != nil {
		log.Fatal(err)
	}
	if !ok {
		log.Printf("Subscription %v doesn't exist. Will create.\n", subName)
		newSub, err := client.CreateSubscription(ctx, subName, pubsub.SubscriptionConfig{
			Topic:            topic,
			AckDeadline:      10 * time.Second,
			ExpirationPolicy: time.Duration(0),
		})
		if err != nil {
			log.Fatal(err)
		}
		sub = newSub
	}
	return
}

func receiveMessages(ctx context.Context, sub *pubsub.Subscription) {
	err := sub.Receive(ctx, func(ctx context.Context, m *pubsub.Message) {
		fmt.Println("Received message:")
		var dat map[string]interface{}
		if err := json.Unmarshal(m.Data, &dat); err != nil {
			log.Println(err)
		}
		fmt.Println(json.MarshalIndent(dat, "", "    "))
		fmt.Println()
		if dat["status"] == "SUCCESS" || dat["status"] == "FAILURE" {
			repo := dat["substitutions"].(map[string]interface{})["REPO_NAME"]
			commit := dat["substitutions"].(map[string]interface{})["COMMIT_SHA"]
			message := fmt.Sprintf("CloudBuild notification: \nStatus: %v\nCommit: https://github.com/raynix/%v/commit/%v\nBuild log: %v\n", dat["status"], repo, commit, dat["logUrl"])
			postToSlack(message)
		}
		m.Ack()
	})
	log.Fatal(err)
}

func postToSlack(text string) {
	payload := slackMessage{Text: text}
	jsonValue, _ := json.Marshal(payload)

	token := os.Getenv("SLACK_TOKEN")
	if len(token) == 0 {
		log.Fatal("SLACK_TOKEN not set!")
	}
	response, _ := http.Post("https://hooks.slack.com/services/"+token, "application/json", bytes.NewBuffer(jsonValue))
	log.Printf("Request sent. Response code is %v: %v\n", response.Status, response.StatusCode)
}

func main() {
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		log.Fatal(err)
	}

	topic := client.Topic(topicName)
	ok, err := topic.Exists(ctx)
	if err != nil {
		log.Fatal(err)
	}
	if !ok {

	}
	sub := findOrCreateSub(ctx, client, topic, subName)
	receiveMessages(ctx, sub)

}
