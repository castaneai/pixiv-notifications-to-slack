package pixiv_notifications_to_slack

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"
)

const (
	baseURL           = "https://www.pixiv.net"
	sessionCookieName = "PHPSESSID"
)

type Notification struct {
	ID         string    `json:"id" firestore:"id"`
	Content    string    `json:"content" firestore:"content"`
	NotifiedAt time.Time `json:"notifiedAt" firestore:"notifiedAt"`
	LinkURL    string    `json:"linkUrl" firestore:"linkUrl"`
	IconURL    string    `json:"iconUrl" firestore:"iconUrl"`
}

// pixiv notification ID's type is (int|string) 😨
type ambiguousNotification struct {
	ID         json.RawMessage `json:"id"`
	Content    string          `json:"content"`
	NotifiedAt time.Time       `json:"notifiedAt"`
	LinkURL    string          `json:"linkUrl"`
	IconURL    string          `json:"iconUrl"`
}

func (b *ambiguousNotification) ToNotification() *Notification {
	// Unquote string in json.RawMessage "\"string\"" -> "string"
	id := strings.Trim(fmt.Sprintf("%s", b.ID), "\"")
	return &Notification{
		ID:         id,
		Content:    b.Content,
		NotifiedAt: b.NotifiedAt,
		LinkURL:    b.LinkURL,
		IconURL:    b.IconURL,
	}
}

type responseJSONBody struct {
	Items []*ambiguousNotification `json:"items"`
}

type responseJSON struct {
	Error   bool              `json:"error"`
	Message string            `json:"message"`
	Body    *responseJSONBody `json:"body"`
}

func (rb *responseJSON) UnmarshalJSON(b []byte) error {
	var a struct {
		Error   bool   `json:"error"`
		Message string `json:"message"`
	}
	if err := json.Unmarshal(b, &a); err != nil {
		return err
	}
	rb.Error = a.Error
	rb.Message = a.Message
	if !a.Error {
		var valid struct {
			Body *responseJSONBody `json:"body"`
		}
		if err := json.Unmarshal(b, &valid); err != nil {
			return err
		}
		rb.Body = valid.Body
	}
	return nil
}

func GetNotifications(ctx context.Context, sessionID string) ([]*Notification, error) {
	req, err := http.NewRequest("GET", baseURL+"/ajax/notification", nil)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	req.AddCookie(createSessionCookie(sessionID))
	hc := &http.Client{}
	res, err := hc.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var rj responseJSON
	if err := json.NewDecoder(res.Body).Decode(&rj); err != nil {
		return nil, err
	}
	if rj.Error {
		return nil, errors.New(rj.Message)
	}

	var result []*Notification
	// reverse slice
	for i := len(rj.Body.Items) - 1; i >= 0; i-- {
		result = append(result, rj.Body.Items[i].ToNotification())
	}
	return result, nil
}

func createSessionCookie(session string) *http.Cookie {
	expires := time.Now().AddDate(1, 0, 0)
	return &http.Cookie{Name: sessionCookieName, Value: session, Expires: expires, HttpOnly: true}
}
