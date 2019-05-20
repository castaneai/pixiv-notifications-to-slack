package pixiv_notifications_to_slack

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"

	"github.com/stretchr/testify/assert"

	firebase "firebase.google.com/go"
)

func cleanUpCollection(ctx context.Context, client *firestore.Client, collectionPath string) error {
	iter := client.Collection(collectionPath).Documents(ctx)
	defer iter.Stop()
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if _, err := doc.Ref.Delete(ctx); err != nil {
			return err
		}
	}
	return nil
}

func newTestFirestoreStore(ctx context.Context) (Store, error) {
	if os.Getenv("FIRESTORE_EMULATOR_HOST") == "" {
		return nil, fmt.Errorf("FIRESTORE_EMULATOR_HOST not set")
	}
	fb, err := firebase.NewApp(ctx, nil)
	if err != nil {
		return nil, err
	}
	client, err := fb.Firestore(ctx)
	if err != nil {
		return nil, err
	}

	collectionPath := "pixivNotifications"
	if err := cleanUpCollection(ctx, client, collectionPath); err != nil {
		return nil, err
	}
	return NewFirestoreStore(client, collectionPath)
}

func Test_MarkAsRead(t *testing.T) {
	ctx := context.Background()
	store, err := newTestFirestoreStore(ctx)
	if err != nil {
		t.Fatal(err)
	}

	firstNotifiedAt := time.Now()
	ns := []*Notification{
		{ID: 1, Content: "n1", NotifiedAt: firstNotifiedAt},
		{ID: 2, Content: "n2", NotifiedAt: firstNotifiedAt},
	}

	{
		unreads, err := store.Unreads(ctx, ns)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, 2, len(unreads))

		if err := store.MarkAsRead(ctx, unreads); err != nil {
			t.Fatal(err)
		}
	}

	{
		unreads, err := store.Unreads(ctx, ns)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, 0, len(unreads))
	}

	{
		secondNotifiedAt := firstNotifiedAt.Add(1 * time.Second)
		ns := append(ns, &Notification{ID: 3, Content: "n3", NotifiedAt: secondNotifiedAt})
		unreads, err := store.Unreads(ctx, ns)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, 1, len(unreads))
		assert.Equal(t, "n3", unreads[0].Content)

		if err := store.MarkAsRead(ctx, unreads); err != nil {
			t.Fatal(err)
		}
	}

	{
		unreads, err := store.Unreads(ctx, ns)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, 0, len(unreads))
	}
}
