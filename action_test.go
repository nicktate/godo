package godo

import (
	"fmt"
	"net/http"
	"testing"
	"time"
)

func TestAction_ActionsServiceOpImplementsActionsService(t *testing.T) {
	if !Implements((*ActionsService)(nil), new(ActionsServiceOp)) {
		t.Error("ActionsServiceOp does not implement ActionsService")
	}
}

func TestAction_List(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/actions", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"actions": [{"id":1},{"id":2}]}`)
		testMethod(t, r, "GET")
	})

	actions, _, err := client.Actions.List()
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	expected := []Action{{ID: 1}, {ID: 2}}
	if len(actions) != len(expected) || actions[0].ID != expected[0].ID || actions[1].ID != expected[1].ID {
		t.Fatalf("unexpected response")
	}
}

func TestAction_Get(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/actions/12345", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"action": {"id":12345}}`)
		testMethod(t, r, "GET")
	})

	action, _, err := client.Actions.Get(12345)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if action.ID != 12345 {
		t.Fatalf("unexpected response")
	}
}

func TestAction_String(t *testing.T) {
	pt, err := time.Parse(time.RFC3339, "2014-05-08T20:36:47Z")
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	startedAt := &Timestamp{
		Time: pt,
	}
	action := &Action{
		ID:        1,
		Status:    "in-progress",
		Type:      "transfer",
		StartedAt: startedAt,
	}

	stringified := action.String()
	expected := `godo.Action{ID:1, Status:"in-progress", Type:"transfer", ` +
		`StartedAt:godo.Timestamp{2014-05-08 20:36:47 +0000 UTC}, ` +
		`ResourceID:0, ResourceType:""}`
	if expected != stringified {
		t.Errorf("Action.Stringify returned %+v, expected %+v", stringified, expected)
	}
}