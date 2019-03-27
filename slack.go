package pixiv_notifications_to_slack

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/grokify/html-strip-tags-go"
	"log"
	"net/http"
)

type SlackMessage struct {
	Text    string `json:"text"`
	IconURL string `json:"icon_url,omitempty"`
}

func postSlackMessage(ctx context.Context, slackWebHookURL string, message *SlackMessage) error {
	buf, err := json.Marshal(message)
	if err != nil {
		return err
	}
	resp, err := http.Post(slackWebHookURL, "application/json", bytes.NewReader(buf))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	return nil
}

func postNotificationToSlack(ctx context.Context, slackWebhookURL string, n *Notification) error {
	text := strip.StripTags(n.Content)
	mes := &SlackMessage{
		Text:    text,
		IconURL: n.IconURL,
	}
	log.Printf("post to slack: %+v", mes)
	err := postSlackMessage(ctx, slackWebhookURL, mes)
	if err != nil {
		return err
	}
	return nil
}
