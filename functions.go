package pixiv_notifications_to_slack

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	firebase "firebase.google.com/go"

	"cloud.google.com/go/firestore"
)

const (
	firestoreCollectionPath = "pixivNotifications"
)

var store *firestore.Client

func PixivNotificationsToSlack(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if store == nil {
		st, err := initStore(ctx)
		if err != nil {
			log.Printf("error initStore: %+v", err)
			w.WriteHeader(500)
			return
		}
		store = st
	}

	if err := check(ctx, store); err != nil {
		log.Printf("error check: %+v", err)
		w.WriteHeader(500)
	}

	w.WriteHeader(200)
}

func initStore(ctx context.Context) (*firestore.Client, error) {
	fa, err := firebase.NewApp(ctx, nil)
	if err != nil {
		return nil, err
	}
	fs, err := fa.Firestore(ctx)
	if err != nil {
		return nil, err
	}
	return fs, err
}

func check(ctx context.Context, client *firestore.Client) error {
	sessionID := os.Getenv("PIXIV_SESSION")
	if sessionID == "" {
		return fmt.Errorf("env not set: PIXIV_SESSION")
	}
	ns, err := GetNotifications(ctx, sessionID)
	if err != nil {
		return err
	}

	slackWebhookURL := os.Getenv("SLACK_WEBHOOK_URL")
	if slackWebhookURL == "" {
		return fmt.Errorf("env not set: SLACK_WEBHOOK_URL")
	}

	store, err := NewFirestoreStore(client, firestoreCollectionPath)
	if err != nil {
		return err
	}
	unreads, err := store.Unreads(ctx, ns)
	if err != nil {
		return err
	}
	if len(unreads) < 1 {
		return nil
	}

	if err := store.MarkAsRead(ctx, unreads); err != nil {
		return err
	}

	for _, n := range unreads {
		if err := postNotificationToSlack(ctx, slackWebhookURL, n); err != nil {
			return err
		}
	}
	return nil
}
