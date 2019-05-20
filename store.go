package pixiv_notifications_to_slack

import (
	"context"
)

type Store interface {
	Unreads(context.Context, []*Notification) ([]*Notification, error)
	MarkAsRead(context.Context, []*Notification) error
}
