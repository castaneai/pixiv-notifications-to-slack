package pixiv_notifications_to_slack

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUnmarshalJSON(t *testing.T) {
	success := []byte(`{"error": false, "message": "success", "body": {"items": []}}`)
	var r responseJSON
	if err := json.Unmarshal(success, &r); err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "success", r.Message)
	assert.Equal(t, 0, len(r.Body.Items))

	failed := []byte(`{"error": true, "message": "error!", "body": []}`)
	var r2 responseJSON
	if err := json.Unmarshal(failed, &r2); err != nil {
		t.Fatal(err)
	}
	assert.Nil(t, r2.Body)
}
