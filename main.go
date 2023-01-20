package main

import (
	"bytes"
	"cloud.google.com/go/pubsub"
	"github.com/slack-go/slack"
	"context"
	"fmt"
	"log"
	"os"
	"text/template"
)

const (
	fallbackSlackTemplate = `[{"type":"section","text":{"type":"plain","text":"{{ printf "%s" .Data }}"}}]`
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
	tmpl, err := template.ParseFiles("slack.json")
	if err != nil {
		tmpl = template.Must(template.New("slack").Parse(fallbackSlackTemplate))
	}
	err = subscribe(projectID, subID, slackWebhookUrl, tmpl)
	if err != nil {
		log.Printf("Receive: %v \n", err)
	}	
}

func subscribe(proj string, subID string, slackWebhook string, tmpl *template.Template) error {
        ctx := context.Background()
	client, err := pubsub.NewClient(ctx, proj)
	if err != nil {
		return fmt.Errorf("pubsub.NewClient: %v", err)
	}
	defer client.Close()
	sub := client.Subscription(subID)
	err = sub.Receive(ctx, func(ctx context.Context, m *pubsub.Message) {
		msg, err := formatMessage(m, tmpl)
		if err != nil {
			m.Nack()
			log.Fatalf("failed to write Slack message: %v", err)
		}
		slack.PostWebhook(slackWebhook, msg)
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

func formatMessage(m *pubsub.Message, tmpl *template.Template) (*slack.WebhookMessage, error) {
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, m); err != nil {
		return nil, err
	}
	var blocks slack.Blocks
	err := blocks.UnmarshalJSON(buf.Bytes())
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal templating JSON: %w", err)
	}

	var clr string
	switch m.Attributes["type_url"] {
	case "type.googleapis.com/google.container.v1beta1.UpgradeAvailableEvent":
		clr = "good"
	case "type.googleapis.com/google.container.v1beta1.UpgradeEvent":
		clr = "warning"
	case "type.googleapis.com/google.container.v1beta1.SecurityBulletinEvent":
		clr = ":exclamation:"
	default:
		clr = "danger"
	}
	return &slack.WebhookMessage{Attachments: []slack.Attachment{{Color: clr}}, Blocks: &blocks}, nil
}
