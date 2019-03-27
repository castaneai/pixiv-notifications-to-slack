package pixiv_notifications_to_slack

import (
	"cloud.google.com/go/firestore"
	"context"
)

const (
	collectionPath = "pixivNotifications"
)

type Store struct {
	client *firestore.Client
}

func (st *Store) FilterUnreadNotifications(ctx context.Context, ns []*Notification) ([]*Notification, error) {
	latestRead, err := st.latestReadNotification(ctx)
	if err != nil {
		return nil, err
	}

	var unreads []*Notification
	for _, n := range ns {
		if isUnread(n, latestRead) {
			unreads = append(unreads, n)
		}
	}

	return unreads, nil
}

func (st *Store) latestReadNotification(ctx context.Context) (*Notification, error) {
	q := st.client.Collection(collectionPath).OrderBy("notifiedAt", firestore.Desc).Limit(1)
	iter := q.Documents(ctx)
	defer iter.Stop()

	var latest Notification
	doc, err := iter.Next()
	if err != nil {
		return nil, err
	}
	if !doc.Exists() {
		return nil, nil
	}
	if err := doc.DataTo(&latest); err != nil {
		return nil, err
	}
	return &latest, nil
}

func (st *Store) MarkAsRead(ctx context.Context, ns []*Notification) error {
	cref := st.client.Collection(collectionPath)
	batch := st.client.Batch()
	for _, n := range ns {
		doc := cref.NewDoc()
		batch = batch.Create(doc, n)
	}
	if _, err := batch.Commit(ctx); err != nil {
		return err
	}
	return nil
}

func isUnread(n, latestRead *Notification) bool {
	if latestRead == nil {
		return true
	}
	return n.ID != latestRead.ID && n.NotifiedAt.After(latestRead.NotifiedAt)
}