package pixiv_notifications_to_slack

import (
	"context"

	"google.golang.org/api/iterator"

	"cloud.google.com/go/firestore"
)

type FirestoreStore struct {
	client         *firestore.Client
	collectionPath string
}

func NewFirestoreStore(client *firestore.Client, collectionPath string) (Store, error) {
	return &FirestoreStore{
		client:         client,
		collectionPath: collectionPath,
	}, nil
}

func (fs *FirestoreStore) Unreads(ctx context.Context, ns []*Notification) ([]*Notification, error) {
	latestRead, err := fs.latestReadNotification(ctx)
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

func (fs *FirestoreStore) latestReadNotification(ctx context.Context) (*Notification, error) {
	q := fs.client.Collection(fs.collectionPath).
		OrderBy("notifiedAt", firestore.Desc).
		Limit(1)

	iter := q.Documents(ctx)
	defer iter.Stop()

	var latest Notification
	doc, err := iter.Next()
	if err == iterator.Done {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	if err := doc.DataTo(&latest); err != nil {
		return nil, err
	}
	return &latest, nil
}

func (fs *FirestoreStore) MarkAsRead(ctx context.Context, ns []*Notification) error {
	cref := fs.client.Collection(fs.collectionPath)
	batch := fs.client.Batch()
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
