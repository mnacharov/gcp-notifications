package main

import (
	"cloud.google.com/go/pubsub"
	"github.com/slack-go/slack"
	"github.com/google/go-cmp/cmp"
	"testing"
	"text/template"
)


func TestWriteMessage(t *testing.T) {
	tmpl, err := template.New("slack").Parse(fallbackSlackTemplate)
	if err != nil {
		t.Fatalf("failed to parse fallback template: %v", err)
	}
	tmpl, err = template.ParseFiles("slack.json")
	if err != nil {
		t.Fatalf("failed to parse default template: %v", err)
	}

	m := &pubsub.Message{
		ID: "ID",
		Data: []byte("qwe qwe"),
		Attributes: map[string]string{
			"type_url": "type.googleapis.com/google.container.v1beta1.UpgradeEvent",
			"project_id": "713594071372",
			"cluster_name": "gke",
			"payload": "{\"resourceType\":\"NODE_POOL\",\"operation\":\"operation-1674134817535-4e475850\",\"operationStartTime\":\"2023-01-19T13:26:57.535908592Z\",\"currentVersion\":\"1.23.13-gke.900\",\"targetVersion\":\"1.23.14-gke.1800\",\"resource\":\"projects/qwe-qwe/locations/europe-west4-a/clusters/gke/nodePools/system\"}",
			"cluster_location": "europe-west4-a",
		},
	}
	got, err := formatMessage(m, tmpl)
	if err != nil {
		t.Fatalf("formatMessage failed: %v", err)
	}

	want := &slack.WebhookMessage{
		Attachments: []slack.Attachment{{Color: "warning"}},
		Blocks: &slack.Blocks{
			BlockSet: []slack.Block{
				&slack.SectionBlock{
					Type: "section",
					Text: &slack.TextBlockObject{
						Type: "mrkdwn",
						Text: "qwe qwe\n`{\"resourceType\":\"NODE_POOL\",\"operation\":\"operation-1674134817535-4e475850\",\"operationStartTime\":\"2023-01-19T13:26:57.535908592Z\",\"currentVersion\":\"1.23.13-gke.900\",\"targetVersion\":\"1.23.14-gke.1800\",\"resource\":\"projects/qwe-qwe/locations/europe-west4-a/clusters/gke/nodePools/system\"}`",
					},
				},
				&slack.DividerBlock{
					Type: "divider",
				},
				&slack.SectionBlock{
					Type: "section",
					Text: &slack.TextBlockObject{
						Type: "mrkdwn",
						Text: "Google Cloud console",
					},
					Accessory: &slack.Accessory{ButtonElement: &slack.ButtonBlockElement{
						Type:     "button",
						Text:     &slack.TextBlockObject{Type: "plain_text", Text: "gke"},
						ActionID: "button-action",
						URL:      "https://console.cloud.google.com/kubernetes/clusters/details/europe-west4-a/gke/details?project=713594071372",
						Value:    "open",
					}},
				},
			},
		},
	}

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("formatMessage got unexpected diff: %s", diff)
	}
}
