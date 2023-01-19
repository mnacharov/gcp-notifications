package main

import (
	"bytes"
	"cloud.google.com/go/pubsub"
	"context"
	"fmt"
	"log"
	"os"
)


func main() {
	projectID := os.Getenv("GCP_PROJECT_ID")
	if projectID == "" {
		log.Fatal("GCP_PROJECT_ID must be set")
		return
	}
	subID := os.Getenv("GCP_SUBSCRIPTION_ID")
	if subID == "" {
		log.Fatal("GCP_SUBSCRIPTION_ID must be set")
		return
	}
	slackWebhookUrl := os.Getenv("SLACK_WEBHOOK_URL")
	if slackWebhookUrl == "" {
		log.Fatal("SLACK_WEBHOOK_URL must be set")
		return
	}
	err := subscribe(projectID, subID, slackWebhookUrl)
	if err != nil {
		log.Printf("Receive: %v \n", err)
	}	
}

func subscribe(proj string, subID string, slack string) error {
        ctx := context.Background()
	client, err := pubsub.NewClient(ctx, proj)
	if err != nil {
		return fmt.Errorf("pubsub.NewClient: %v", err)
	}
	defer client.Close()
	sub := client.Subscription(subID)
	err = sub.Receive(ctx, func(ctx context.Context, m *pubsub.Message) {
		msg := formatMessage(m)
		sendSlack(slack, msg)
		if err == nil {
			m.Ack()
		} else {
			m.Nack()
		}
	})
	if err != context.Canceled {
		return fmt.Errorf("sub.Receive: %v", err)
	}
	return nil
}

func formatMessage(m *pubsub.Message) string {
	var buff bytes.Buffer
	buff.Write(m.Data)
	for k, v := range m.Attributes {
		buff.WriteString(fmt.Sprintf("\n\t%s:\t%s", k, v))
	}
	return buff.String()
}
