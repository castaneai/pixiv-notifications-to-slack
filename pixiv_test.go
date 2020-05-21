package pixiv_notifications_to_slack

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnmarshalJSON(t *testing.T) {
	success := []byte(`{"error": false, "message": "", "body": {"items": [
	{"id": "string-id", "content": "test"},
	{"id": 12345, "content": "test"}
]}}`)
	var r responseJSON
	if err := json.Unmarshal(success, &r); err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "", r.Message)
	assert.Equal(t, 2, len(r.Body.Items))
	n1 := r.Body.Items[0].ToNotification()
	assert.Equal(t, "string-id", n1.ID)
	n2 := r.Body.Items[1].ToNotification()
	assert.Equal(t, "12345", n2.ID)

	failed := []byte(`{"error": true, "message": "error!", "body": []}`)
	var r2 responseJSON
	if err := json.Unmarshal(failed, &r2); err != nil {
		t.Fatal(err)
	}
	assert.Nil(t, r2.Body)
}
