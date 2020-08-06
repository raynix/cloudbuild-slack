package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"cloud.google.com/go/pubsub"
)

var projectID = "idyllic-depth-239301"
var topicName = "cloud-builds"
var subName = "cloud-build-slack-sub"

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
		fmt.Printf("------\n%v\n------\n", m.Data)
		m.Ack()
	})
	log.Fatal(err)
}

func main() {
	fmt.Printf("hello, world\n")
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
